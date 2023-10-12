package service

import (
	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
	"go.jhphx.com/btmn-backend/array"
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

func (s LinkService) GetTags() ([]string, error) {
	links, err := s.GetLinks()
	if err != nil {
		return nil, err
	}

	var tags []string
	for _, l := range links {
		for _, t := range l.Tags {
			if !array.Contains(tags, t) {
				tags = append(tags, t)
			}
		}
	}
	return tags, nil
}
