package service

import (
	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
)

type LinkService struct {
	linkRepository port.LinkRepository
}

func NewLinkService(linkRepository port.LinkRepository) *LinkService {
	return &LinkService{linkRepository: linkRepository}

}

func (s LinkService) Create(link domain.Link) (domain.Link, error) {
	return s.linkRepository.Create(link)
}

func (s LinkService) GetLinks() ([]domain.Link, error) {
	return s.linkRepository.GetLinks()
}
