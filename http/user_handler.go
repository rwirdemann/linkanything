package http

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/linkanything/core"
	"io"
	"log"
	"net/http"
	"strings"
)

type UserHTTPHandler struct {
	service *core.UserService
}

func NewUserHTTPHandler(service *core.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{service: service}
}

func (h UserHTTPHandler) Create() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)
		b, err := io.ReadAll(request.Body)
		if err != nil || len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		var user core.User
		err = json.Unmarshal(b, &user)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		u, err := h.service.Create(user)
		if err != nil {
			log.Print(err)
			if strings.HasPrefix(err.Error(), "user exists") {
				writer.WriteHeader(http.StatusConflict)
			}
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := request.URL.String()
		writer.Header().Set("Location", fmt.Sprintf("%s/%d", url, u.Id))
		writer.WriteHeader(http.StatusCreated)
	}
}
