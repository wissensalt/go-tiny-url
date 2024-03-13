package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/wissensalt/go-tiny-url/config"
	http2 "github.com/wissensalt/go-tiny-url/internal/http"
	"github.com/wissensalt/go-tiny-url/internal/repository"
	"github.com/wissensalt/go-tiny-url/internal/service"
	"net/http"
)

var sqlDB *sql.DB

func main() {
	sqlDB = config.ConnectDB()
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Golang Tiny URL Web Service"))
	})

	r.Mount("/tiny-urls", TinyUrlRoutes())

	_ = http.ListenAndServe("localhost:8080", r)
}

func TinyUrlRoutes() chi.Router {
	urlRepository := repository.UrlRepositoryImpl{DB: sqlDB}
	urlService := service.UrlServiceImpl{UrlRepositoryImpl: urlRepository}
	urlController := http2.UrlControllerImpl{UrlServiceImpl: urlService}
	r := chi.NewRouter()
	r.Get("/", urlController.GetUrls)
	r.Post("/", urlController.Shorten)
	r.Get("/{code}", urlController.Redirect)

	return r
}
