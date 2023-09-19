package ports

import "github.com/rwirdemann/linkanything/domain"

type LinkRepository interface {
	Create(link domain.Link) domain.Link
}
