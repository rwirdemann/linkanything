package linkanything

type LinkRepository interface {
	Create(link Link) (Link, error)
	GetLinks(tags []string, includeDrafts bool, page, limit int) ([]Link, error)
	Get(id int) (Link, error)
	Update(l Link) (Link, error)
	Count() (int, error)
}
