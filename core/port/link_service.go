package port

import "github.com/rwirdemann/linkanything/core/domain"

type LinkService interface {
	Create(link domain.Link) domain.Link
}
