package yaml

import (
	"time"

	r "github.com/freerware/tutor/api/representations"
	"github.com/freerware/tutor/domain"
	u "github.com/gofrs/uuid"
)

type Post struct {
	r.Representation `yaml:"-"`

	UUID      u.UUID     `yaml:"uuid"`
	Title     string     `yaml:"title"`
	Content   string     `yaml:"content"`
	Draft     bool       `yaml:"isDraft"`
	Likes     int        `yaml:"likes"`
	CreatedAt time.Time  `yaml:"createdAt"`
	UpdatedAt time.Time  `yaml:"updatedAt"`
	DeletedAt *time.Time `yaml:"deletedAt"`
}

// Bytes provides the representation as bytes.
func (p Post) Bytes() ([]byte, error) {
	return p.Base.Bytes(&p)
}

// FromBytes constructs the representation from bytes.
func (p Post) FromBytes(b []byte) error {
	return p.Base.FromBytes(b, &p)
}

// NewPost constructs a new account representation.
func NewPost(p domain.Post) Post {
	post := Post{
		UUID:      p.UUID(),
		Title:     p.Title(),
		Content:   p.Content(),
		Draft:     p.IsDraft(),
		Likes:     p.Likes(),
		CreatedAt: p.CreatedAt(),
		UpdatedAt: p.UpdatedAt(),
		DeletedAt: p.DeletedAt(),
	}
	post.SetContentCharset("ascii")
	post.SetContentLanguage("en-US")
	post.SetContentType("application/yaml")
	post.SetSourceQuality(1.0)
	post.SetContentEncoding([]string{"identity"})
	return post
}

func NewPosts(posts ...domain.Post) []Post {
	postRepresentations := make([]Post, len(posts))
	for i, post := range posts {
		postRepresentations[i] = NewPost(post)
	}
	return postRepresentations
}
