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
	repository port.LinkRepository
}

func NewHTTPHandler(service port.LinkService) *HTTPHandler {
	return &HTTPHandler{repository: service}
}

func (h HTTPHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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