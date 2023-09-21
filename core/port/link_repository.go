package port

import "github.com/rwirdemann/linkanything/core/domain"

type LinkRepository interface {
	Create(link domain.Link) (domain.Link, error)
	GetLinks() ([]domain.Link, error)
}
