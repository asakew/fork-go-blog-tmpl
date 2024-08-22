package models

import "time"

type Article struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	PublishedDate time.Time `json:"published_date"`
}

type ArticleResponse struct {
	Article *Article `json:"article"`
}

type ArticlesResponse struct {
	Articles []Article `json:"articles"`
}
