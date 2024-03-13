package repository

import (
	"database/sql"
	"errors"
	"math/rand"
)

type (
	Url struct {
		Id        int
		Code      string
		OriginUrl string
	}

	UrlRepository interface {
		GetUrls() []Url
		Shorten(originUrl string) (Url, error)
		FindByCode(code string) (Url, error)
	}

	UrlRepositoryImpl struct {
		*sql.DB
	}
)

func (u UrlRepositoryImpl) GetUrls() []Url {
	query := "SELECT * FROM url"
	rows, err := u.DB.Query(query)
	if err != nil {
		return []Url{}
	}

	var urls = make([]Url, 0)
	for rows.Next() {
		var url Url
		err := rows.Scan(&url.Id, &url.Code, &url.OriginUrl)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return urls
			}
		}

		urls = append(urls, url)
	}

	return urls
}

func (u UrlRepositoryImpl) Shorten(originUrl string) (Url, error) {
	query := "SELECT * FROM url WHERE origin_url=$1"
	r := u.DB.QueryRow(query, originUrl)
	var url Url
	err := r.Scan(&url.Id, &url.Code, &url.OriginUrl)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return url, err
	} else {
		code := generateRandomString(7)
		query = "INSERT INTO url(code, origin_url) VALUES ($1, $2)"
		_, err := u.DB.Exec(query, code, originUrl)
		if err != nil {
			return url, err
		}

		return Url{
			Id:        0,
			Code:      code,
			OriginUrl: originUrl,
		}, nil
	}
}

func (u UrlRepositoryImpl) FindByCode(code string) (Url, error) {
	r := u.DB.QueryRow("SELECT * FROM url WHERE code=$1", code)
	var url Url
	err := r.Scan(&url.Id, &url.Code, &url.OriginUrl)
	if err != nil {
		return Url{}, err
	}

	return url, nil
}

// https://stackoverflow.com/a/22892986/23602756
func generateRandomString(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	r := make([]rune, n)
	for i := range r {
		r[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(r)
}
