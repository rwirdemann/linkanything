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

func (s LinkService) GetLinks(tags []string) ([]domain.Link, error) {
	return s.linkRepository.GetLinks(tags)
}

func (s LinkService) GetTags() ([]string, error) {
	links, err := s.GetLinks(nil)
	if err != nil {
		return nil, err
	}

	var tags []string
	for _, l := range links {
		for _, t := range l.Tags {
			if !contains(tags, t) {
				tags = append(tags, t)
			}
		}
	}
	return tags, nil
}

func contains(a []string, e string) bool {
	for _, v := range a {
		if e == v {
			return true
		}
	}
	return false
}
