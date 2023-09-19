package domain

import "time"

type Link struct {
	Id      int       `json:"id"`
	Title   string    `json:"title"`
	Type    string    `json:"type"`
	Body    string    `json:"body"`
	URI     string    `json:"uri"`
	Created time.Time `json:"created"`
	Tags    []string  `json:"tags"`
}
