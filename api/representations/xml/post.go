package xml

import (
	"time"

	r "github.com/freerware/tutor/api/representations"
	"github.com/freerware/tutor/domain"
	u "github.com/gofrs/uuid"
)

type Post struct {
	r.Representation `xml:"-"`

	UUID      u.UUID     `xml:"uuid"`
	Title     string     `xml:"title"`
	Content   string     `xml:"content"`
	Draft     bool       `xml:"isDraft"`
	Likes     int        `xml:"likes"`
	CreatedAt time.Time  `xml:"createdAt"`
	UpdatedAt time.Time  `xml:"updatedAt"`
	DeletedAt *time.Time `xml:"deletedAt"`
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
	post.SetContentType("application/xml")
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
