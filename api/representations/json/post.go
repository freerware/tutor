package json

import (
	"time"

	r "github.com/freerware/tutor/api/representations"
	"github.com/freerware/tutor/domain"
	u "github.com/gofrs/uuid"
)

type Post struct {
	r.Representation `json:"-"`

	UUID      u.UUID     `json:"uuid"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Draft     bool       `json:"isDraft"`
	Likes     int        `json:"likes"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
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
	post.SetContentType("application/json")
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
