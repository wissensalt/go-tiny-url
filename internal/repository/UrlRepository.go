package repository

import (
	"database/sql"
	"errors"
	"fmt"
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
	query := "INSERT INTO url(code, origin_url) VALUES ($1, $2) RETURNING id"
	tx, err := u.DB.Begin()
	if err != nil {
		fmt.Println("Failed to start transaction")
		return Url{}, err
	}

	stmt, err := tx.Prepare(query)
	if err != nil {
		fmt.Println("Failed to create statement", err.Error())
		return Url{}, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			fmt.Println("Failed to close statement", err.Error())
		}
	}(stmt)

	var savedUrl Url
	var urlId int
	generatedCode := generateRandomString(7)
	err = stmt.QueryRow(generatedCode, originUrl).Scan(&urlId)
	if err != nil {
		fmt.Println("Failed to get last inserted id", err.Error())
		return Url{}, err
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Failed to commit transaction")
		return Url{}, err
	}

	savedUrl.Id = urlId
	savedUrl.Code = generatedCode
	savedUrl.OriginUrl = originUrl

	return savedUrl, nil
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
