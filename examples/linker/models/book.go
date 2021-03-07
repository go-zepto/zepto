package models

import "time"

type Book struct {
	Model
	Name        string     `json:"name"`
	AuthorID    *uint      `json:"author_id"`
	Author      *Author    `json:"author"`
	PublishedAt *time.Time `json:"published_at"`
}
