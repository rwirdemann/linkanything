package adapter

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/linkanything/domain"
	"github.com/rwirdemann/linkanything/ports"
	"io"
	"net/http"
)

type LinkHandler struct {
	repository ports.LinkRepository
}

func NewLinkHandler(repository ports.LinkRepository) *LinkHandler {
	return &LinkHandler{repository: repository}
}

func (h LinkHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		b, err := io.ReadAll(request.Body)
		if err != nil || len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
		}
		var link domain.Link
		err = json.Unmarshal(b, &link)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
		}

		n := h.repository.Create(link)
		url := request.URL.String()
		writer.Header().Set("Location", fmt.Sprintf("%s/%d", url, n.Id))
		writer.WriteHeader(http.StatusCreated)
	}
}
