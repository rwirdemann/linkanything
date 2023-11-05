package port

import "github.com/rwirdemann/linkanything/core/domain"

type LinkService interface {
	Create(link domain.Link) (domain.Link, error)
	GetLinks(tags []string) ([]domain.Link, error)
	GetTags() ([]string, error)
	Get(id int) (domain.Link, error)
}
