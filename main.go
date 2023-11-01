package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	http2 "github.com/rwirdemann/linkanything/adapter/http"
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

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	sessionAdapter := http2.NewSessionHandler()

	linkRepository := adapter.NewPostgresLinkRepository(dbpool)
	linkService := service.NewLinkService(linkRepository)
	linkAdapter := http2.NewLinkHandler(linkService)

	userRepoitory := adapter.NewPostgresUserRepository(dbpool)
	userService := service.NewUserService(userRepoitory)
	userAdapter := http2.NewUserHTTPHandler(userService)

	router := mux.NewRouter()
	router.HandleFunc("/users", userAdapter.Create()).Methods("POST")
	router.HandleFunc("/sessions", sessionAdapter.Create()).Methods("POST")
	router.HandleFunc("/links", linkAdapter.Create()).Methods("POST")
	router.HandleFunc("/links", linkAdapter.GetLinks()).Methods("GET")
	router.HandleFunc("/tags", linkAdapter.GetTags()).Methods("GET")
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

	tls := flag.Bool("tls", true, "use tls")
	flag.Parse()
	if *tls {
		err = http.ListenAndServeTLS(
			fmt.Sprintf(":%s", os.Getenv("PORT")),
			os.Getenv("SSH_PUBLIC_KEY"),
			os.Getenv("SSH_PRIVATE_KEY"),
			router)
	} else {
		err = http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router)
	}

	if err != nil {
		log.Fatal(err)
	}
}
