package http

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/linkanything/core"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type LinkHandler struct {
	repository core.LinkRepository
}

func NewLinkHandler(repository core.LinkRepository) *LinkHandler {
	return &LinkHandler{repository: repository}
}

func (h LinkHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		b, err := io.ReadAll(request.Body)
		if err != nil || len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var link core.Link
		err = json.Unmarshal(b, &link)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		l, err := h.repository.Create(link)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := request.URL.String()
		writer.Header().Set("Location", fmt.Sprintf("%s/%d", url, l.Id))
		writer.WriteHeader(http.StatusCreated)
	}
}

type pagination struct {
	TotalRecords int `json:"total_record"`
	CurrentPage  int `json:"current_page"`
	TotalPages   int `json:"total_pages"`
	NextPage     int `json:"next_page"`
	PrevPage     int `json:"prev_page"`
}

func (h LinkHandler) GetLinks() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)

		// extract tags from request
		tags := request.URL.Query().Get("tags")
		var tagList []string
		if len(tags) > 0 {
			tagList = trim(strings.Split(tags, ","))
		}

		var err error
		includeDrafts := false
		drafts := request.URL.Query().Get("drafts")
		if len(drafts) > 0 {
			includeDrafts, err = strconv.ParseBool(drafts)
			if err != nil {
				log.Fatal(err)
			}
		}

		// extract page from query params
		page := 0
		pageParam := request.URL.Query().Get("page")
		if len(pageParam) > 0 {
			if p, err := strconv.Atoi(pageParam); err == nil {
				page = p
			}
		}

		// extract limit from query params
		limit := 0
		limitParam := request.URL.Query().Get("limit")
		if len(limitParam) > 0 {
			if l, err := strconv.Atoi(limitParam); err == nil {
				limit = l
			}
		}

		links, err := h.repository.GetLinks(tagList, includeDrafts, page, limit)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(
			struct {
				Links      []core.Link `json:"links"`
				Pagination pagination  `json:"pagination"`
			}{
				Links: links,
				Pagination: pagination{
					TotalRecords: len(links),
					CurrentPage:  page,
					TotalPages:   0,
					NextPage:     page + 1,
					PrevPage:     0,
				},
			},
		)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(b)
	}
}

func trim(tags []string) []string {
	var result []string
	for _, t := range tags {
		result = append(result, strings.TrimSpace(t))
	}
	return result
}
