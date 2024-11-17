package http

import (
	"encoding/json"
	"github.com/rwirdemann/linkanything/core"
	"github.com/rwirdemann/linkanything/postgres"
	"io"
	"log"
	"net/http"
)

type SessionHandler struct {
	repository core.UserRepository
}

func NewSessionHandler(repository core.UserRepository) *SessionHandler {
	return &SessionHandler{repository: repository}
}

func (h SessionHandler) Create() func(http.ResponseWriter, *http.Request) {
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
			log.Print(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		hash, _ := h.repository.GetHash(user.Name)
		if err != nil {
			log.Print(err)
			writer.WriteHeader(http.StatusForbidden)
			return
		}

		if !postgres.CheckPasswordHash(user.Password, hash) {
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
