package adapter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
)

type HTTPHandler struct {
	service port.LinkService
}

func NewHTTPHandler(service port.LinkService) *HTTPHandler {
	return &HTTPHandler{service: service}
}

func (h HTTPHandler) Create() http.HandlerFunc {
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

func (h HTTPHandler) GetLinks() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		links, err := h.service.GetLinks()
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

func (h HTTPHandler) GetTags() func(http.ResponseWriter, *http.Request) {
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

func (h HTTPHandler) Login() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
