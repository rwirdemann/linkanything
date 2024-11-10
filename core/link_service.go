package core

import (
	"strconv"
)

type LinkService struct {
	linkRepository LinkRepository
}

func NewLinkService(linkRepository LinkRepository) *LinkService {
	return &LinkService{linkRepository: linkRepository}

}

func (s LinkService) Create(link Link) (Link, error) {
	return s.linkRepository.Create(link)
}

func (s LinkService) GetLinks(tags []string, includeDrafts bool, page, limit int) ([]Link, error) {
	return s.linkRepository.GetLinks(tags, includeDrafts, page, limit)
}

func (s LinkService) GetTags() ([]string, error) {
	links, err := s.GetLinks(nil, true, 0, 0)
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

func (s LinkService) Patch(patch Patch) error {
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

func (s LinkService) Get(id int) (Link, error) {
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
