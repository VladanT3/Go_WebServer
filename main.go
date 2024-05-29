package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/VladanT3/Go_WebServer/internal/database"
	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

    _ "github.com/lib/pq"
)

type apiConfig struct {
    DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	var port string = os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found in environment")
	}

 	var connStr string = os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("DB connection string not found in environment")
	}   

    dbConn, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Couldn't connect to database")
    }

    apiCfg := apiConfig{
        DB: database.New(dbConn),
    }

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)
    v1Router.Post("/users", apiCfg.handlerCreateUser)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server running on: http://localhost:%v\n", port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
