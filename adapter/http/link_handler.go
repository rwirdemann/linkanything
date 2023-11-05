package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
)

type LinkHandler struct {
	service port.LinkService
}

func NewLinkHandler(service port.LinkService) *LinkHandler {
	return &LinkHandler{service: service}
}

func (h LinkHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		b, err := io.ReadAll(request.Body)
		if err != nil || len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var link domain.Link
		err = json.Unmarshal(b, &link)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		l, err := h.service.Create(link)
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

func (h LinkHandler) GetLinks() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)

		// extract tags from request
		tags := request.URL.Query().Get("tags")
		var tagList []string
		if len(tags) > 0 {
			tagList = trim(strings.Split(tags, ","))
		}

		links, err := h.service.GetLinks(tagList)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(
			struct {
				Links []domain.Link `json:"links"`
			}{
				links,
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

func (h LinkHandler) GetTags() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		tags, err := h.service.GetTags()
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(
			struct {
				Tags []string `json:"tags"`
			}{
				tags,
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
