package http

import (
	"encoding/json"
	"github.com/rwirdemann/linkanything/adapter"
	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
	"io"
	"log"
	"net/http"
)

type SessionHandler struct {
	service port.UserService
}

func NewSessionHandler(service port.UserService) *SessionHandler {
	return &SessionHandler{service: service}
}

func (h SessionHandler) Create() func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		enableCors(&writer)

		b, err := io.ReadAll(request.Body)
		if err != nil || len(b) == 0 {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		var user domain.User
		err = json.Unmarshal(b, &user)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		hash, _ := h.service.GetHash(user.Name)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		if !adapter.CheckPasswordHash(user.Password, hash) {
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		token, err := generateJWT(user.Name)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		resp, err := json.Marshal(
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
		writer.Write(resp)
	}
}
