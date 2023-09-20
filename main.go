package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rwirdemann/linkanything/adapter"
	"github.com/rwirdemann/linkanything/core/service"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	linkRepository := adapter.NewPostgresRepository()
	linkService := service.NewLinkService(linkRepository)
	linkAdapter := adapter.NewHTTPHandler(linkService)

	router := mux.NewRouter()
	router.HandleFunc("/links", linkAdapter.Create()).Methods("POST")
	router.HandleFunc("/links", linkAdapter.GetLinks()).Methods("GET")
	log.Printf("starting server on port %s", os.Getenv("PORT"))
	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		log.Println(tpl, met)
		return nil
	})
	if err != nil {
		panic(err)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}
