package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	http3 "github.com/rwirdemann/linkanything/http"
	"github.com/rwirdemann/linkanything/postgres"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	linkRepository := postgres.NewPostgresLinkRepository(dbpool)
	linkAdapter := http3.NewLinkHandler(linkRepository)

	userRepository := postgres.NewPostgresUserRepository(dbpool)
	userAdapter := http3.NewUserHTTPHandler(userRepository)
	sessionAdapter := http3.NewSessionHandler(userRepository)

	router := mux.NewRouter()
	router.HandleFunc("/users", userAdapter.Create()).Methods("POST")
	router.HandleFunc("/sessions", sessionAdapter.Create()).Methods("POST")
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
