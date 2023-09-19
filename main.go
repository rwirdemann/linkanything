package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rwirdemann/linkanything/adapter"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	linkRepository := adapter.NewLinkRepository()
	linkAdapter := adapter.NewLinkHandler(linkRepository)
	router := mux.NewRouter()
	router.HandleFunc("/links", linkAdapter.Create()).Methods("POST")
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
