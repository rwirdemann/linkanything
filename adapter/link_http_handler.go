package adapter

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

type LinkHTTPHandler struct {
	service port.LinkService
}

func NewLinkHTTPHandler(service port.LinkService) *LinkHTTPHandler {
	return &LinkHTTPHandler{service: service}
}

func (h LinkHTTPHandler) Create() http.HandlerFunc {
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

		link.Draft = false
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

func (h LinkHTTPHandler) GetLinks() http.HandlerFunc {
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

func (h LinkHTTPHandler) GetTags() func(http.ResponseWriter, *http.Request) {
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

func (h LinkHTTPHandler) Login() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		token, err := generateJWT()
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(
			struct {
				Token string `json:"token"`
			}{
				token,
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

func (h LinkHTTPHandler) Logout() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		writer.WriteHeader(http.StatusNoContent)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
