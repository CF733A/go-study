package entities

type PetStatus string

const (
	StatusAvailable PetStatus = "available"
	StatusPending   PetStatus = "pending"
	StatusSold      PetStatus = "sold"
)

// Category represents pet category entity
// @Description Business entity for pet category classification
type Category struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"Dogs"`
}

// Tag represents pet tag entity
// @Description Business entity for pet labeling and categorization  
type Tag struct {
	ID   int64  `json:"id" example:"1"`
	Name string `json:"name" example:"friendly"`
}


type Pet struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Category  Category  `json:"category"`
	PhotoURLs []string  `json:"photoUrls"`
	Tags      []Tag     `json:"tags"`
	Status    PetStatus `json:"status"`
}

func (p *Pet) NewPet(name string, category Category, photos []string, tags []Tag) *Pet {
	return &Pet{
		Name:      name,
		Category:  category,
		PhotoURLs: photos,
		Tags:      tags,
		Status:    StatusAvailable,
	}
}

func (p *Pet) MarkAsPending() {
	p.Status = StatusPending
}

func (p *Pet) MarkAsSold() {
	p.Status = StatusSold
}

func (p *Pet) IsAvailable() bool {
	return p.Status == StatusAvailable
}

func (p *Pet) AddPhotoURL(url string) {
	p.PhotoURLs = append(p.PhotoURLs, url)
}

func (p *Pet) AddTag(name string) {
	maxID := int64(0)
	for _, tag := range p.Tags {
		if tag.ID > maxID {
			maxID = tag.ID
		}
	}

	newID := maxID + 1
	p.Tags = append(p.Tags, Tag{ID: newID, Name: name})
}
