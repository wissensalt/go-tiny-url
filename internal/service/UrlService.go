package service

import "github.com/wissensalt/go-tiny-url/internal/repository"

type (
	UrlService interface {
		GetUrls() []repository.Url
		Shorten(originUrl string) repository.Url
		FindByCode(code string) repository.Url
	}

	UrlServiceImpl struct {
		repository.UrlRepositoryImpl
	}
)

func (u UrlServiceImpl) Shorten(originUrl string) repository.Url {
	url, err := u.UrlRepositoryImpl.Shorten(originUrl)
	if err != nil {
		return repository.Url{}
	}

	return url
}

func (u UrlServiceImpl) FindByCode(code string) repository.Url {
	url, err := u.UrlRepositoryImpl.FindByCode(code)
	if err != nil || url == (repository.Url{}) {
		return repository.Url{}
	}

	return url
}
