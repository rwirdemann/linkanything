package service

import (
	"time"

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
	links := []domain.Link{
		{Id: 1, Title: "Cold Hawaii Games 20233", URI: "https://coldhawaiigames.com", Created: time.Now()},
		{Id: 2, Title: "Test: Vayu Aura 2", URI: "https://gleiten.tv/index.php/video/action/view/v/4326/page/530/", Created: time.Now()},
	}
	return links, nil
}
