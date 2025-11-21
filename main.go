package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/harshit18-09/RSS-Aggregator/internal/db"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *db.Queries
}

func main() {
	// feed, err := urlToFeed("https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml") //just for example to see if rss fetching is working
	// if err != nil {
	// 	log.Fatal("cannot get feed:", err)
	// }
	// fmt.Println("Feed Title:", feed.Channel.Title)

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not in env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL not in env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	apiCfg := &apiConfig{
		DB: db.New(conn),
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
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.With(apiCfg.MiddlewareAuth).Get("/users", apiCfg.handlerGetUser)
	v1Router.With(apiCfg.MiddlewareAuth).Post("/feeds", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userContextKey).(db.User)
		apiCfg.handlerCreateFeed(w, r, user)
	})
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.With(apiCfg.MiddlewareAuth).Post("/feed_follows", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userContextKey).(db.User)
		apiCfg.handlerCreateFeedFollow(w, r, user)
	})
	v1Router.With(apiCfg.MiddlewareAuth).Get("/feed_follows", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userContextKey).(db.User)
		apiCfg.handlerGetFeedFollows(w, r, user)
	})
	v1Router.With(apiCfg.MiddlewareAuth).Delete("/feed_follows/{feedFollowID}", func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(userContextKey).(db.User)
		apiCfg.handlerDeleteFeedFollows(w, r, user)
	})
	v1Router.With(apiCfg.MiddlewareAuth).Get("/posts", apiCfg.handlerGetPostsForUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// start scraper in background and it will wait until all comes
	concurrency := 5
	if v := os.Getenv("SCRAPER_CONCURRENCY"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			concurrency = n
		}
	}
	intervalMinutes := 10
	if v := os.Getenv("SCRAPER_INTERVAL_MINUTES"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			intervalMinutes = n
		}
	}
	go startscraping(apiCfg.DB, concurrency, time.Duration(intervalMinutes)*time.Minute)

	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)
}

//json rest api used
