package service

import (
	"github.com/rwirdemann/linkanything/core/domain"
	"github.com/rwirdemann/linkanything/core/port"
	"strconv"
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

func (s LinkService) Patch(patch domain.Patch) error {
	link, err := s.linkRepository.Get(patch.Id)
	if err != nil {
		return err
	}

	if patch.Field == "draft" {
		b, err := strconv.ParseBool(patch.Value)
		if err != nil {
			return err
		}
		link.Draft = b
	}

	_, err = s.linkRepository.Update(link)
	if err != nil {
		return err
	}

	return nil
}

func (s LinkService) Get(id int) (domain.Link, error) {
	return s.linkRepository.Get(id)
}

func contains(a []string, e string) bool {
	for _, v := range a {
		if e == v {
			return true
		}
	}
	return false
}
