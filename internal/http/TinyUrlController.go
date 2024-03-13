package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/wissensalt/go-tiny-url/internal/repository"
	"github.com/wissensalt/go-tiny-url/internal/service"
	"net/http"
)

type (
	UrlRequest struct {
		OriginUrl string
	}

	UrlResponse struct {
		NewUrl string
	}

	UrlController interface {
		GetUrls(writer http.ResponseWriter, request *http.Request)
		Shorten(writer http.ResponseWriter, request *http.Request)
		Redirect(writer http.ResponseWriter, request *http.Request)
	}

	UrlControllerImpl struct {
		service.UrlServiceImpl
	}
)

func (u UrlControllerImpl) GetUrls(writer http.ResponseWriter, request *http.Request) {
	setContentTypeAsJson(writer)
	urls := u.UrlServiceImpl.GetUrls()
	err := json.NewEncoder(writer).Encode(urls)
	if err != nil {
		http.Error(writer, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

func (u UrlControllerImpl) Shorten(writer http.ResponseWriter, request *http.Request) {
	setContentTypeAsJson(writer)
	var urlRequest UrlRequest
	err := json.NewDecoder(request.Body).Decode(&urlRequest)
	if err != nil {
		http.Error(writer, "Failed to parse request body", http.StatusBadRequest)
	}

	url := u.UrlServiceImpl.Shorten(urlRequest.OriginUrl)
	if url == (repository.Url{}) {
		http.Error(writer, "Failed to shorten URL", http.StatusInternalServerError)
	}
	newUrl := "http://localhost:8080/tiny-urls/" + url.Code
	urlResponse := UrlResponse{NewUrl: newUrl}
	err = json.NewEncoder(writer).Encode(urlResponse)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (u UrlControllerImpl) Redirect(writer http.ResponseWriter, request *http.Request) {
	code := chi.URLParam(request, "code")
	url := u.UrlServiceImpl.FindByCode(code)
	if url == (repository.Url{}) {
		http.Error(writer, "Not Found", http.StatusNotFound)
	}

	http.Redirect(writer, request, url.OriginUrl, http.StatusMovedPermanently)
}

func setContentTypeAsJson(writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
}
