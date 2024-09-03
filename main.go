package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elishambadi/go-rss-agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	feed, err := urlToFeed("https://wagslane.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(feed)

	godotenv.Load()

	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	if portString == "" {
		log.Fatal("Port is not found in the env")
	} else {
		fmt.Printf("Port: %s\n", portString)
	}

	if dbUrl == "" {
		log.Fatal("DB URL is not found in the env")
	} else {
		fmt.Printf("db url: %s\n", dbUrl)
	}

	// connect to Db sql.Open("{driver}", "{connection_string}")
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	db := database.New(conn)

	apiCfg := apiConfig{
		DB: db,
	}

	// Call the scraping function to run in the background

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()

	// Cors are used to control access across different clients
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
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeeds))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollowsByUser))
	v1Router.Delete("/feed-follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	v1Router.Delete("/feeds/{feedId}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeed))

	router.Mount("/v1", v1Router)

	// Defining the server
	// As pointer because we don't expect it to change
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err1 := srv.ListenAndServe()
	if err1 != nil {
		log.Fatal(err1)
	}
}
