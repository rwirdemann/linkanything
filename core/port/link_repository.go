package port

import "github.com/rwirdemann/linkanything/core/domain"

type LinkRepository interface {
	Create(link domain.Link) (domain.Link, error)
	GetLinks(tags []string, includeDrafts bool, page, limit int) ([]domain.Link, error)
	Get(id int) (domain.Link, error)
	Update(l domain.Link) (domain.Link, error)
}
