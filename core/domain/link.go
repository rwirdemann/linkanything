package domain

import "time"

type Link struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	URI     string    `json:"uri"`
	Tags    []string  `json:"tags,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Draft   bool
}
