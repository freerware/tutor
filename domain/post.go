package domain

import (
	"time"

	u "github.com/gofrs/uuid"
)

type Post struct {
	uuid       u.UUID
	title      string
	content    string
	draft      bool
	likes      int
	authorUUID u.UUID
	createdAt  time.Time
	updatedAt  time.Time
	deletedAt  *time.Time
}

type PostParameters struct {
	UUID       u.UUID
	Title      string
	Content    string
	Draft      bool
	Likes      int
	AuthorUUID u.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

func NewPost(parameters PostParameters) (Post, error) {
	post := Post{}
	post.SetUUID(parameters.UUID)
	post.SetAuthorUUID(parameters.AuthorUUID)
	post.SetTitle(parameters.Title)
	post.SetContent(parameters.Content)
	post.SetDraft(parameters.Draft)
	post.SetLikes(parameters.Likes)
	if err := post.SetCreatedAt(parameters.CreatedAt); err != nil {
		return Post{}, err
	}
	if err := post.SetUpdatedAt(parameters.UpdatedAt); err != nil {
		return Post{}, err
	}
	if parameters.DeletedAt != nil {
		if err := post.SetDeletedAt(*parameters.DeletedAt); err != nil {
			return Post{}, err
		}
	}
	return post, nil
}

func ReconstitutePost(parameters PostParameters) Post {
	return Post{
		uuid:       parameters.UUID,
		title:      parameters.Title,
		content:    parameters.Content,
		draft:      parameters.Draft,
		likes:      parameters.Likes,
		authorUUID: parameters.AuthorUUID,
		createdAt:  parameters.CreatedAt,
		updatedAt:  parameters.UpdatedAt,
		deletedAt:  parameters.DeletedAt,
	}
}

func (p Post) UUID() u.UUID {
	return p.uuid
}

func (p *Post) SetUUID(uuid u.UUID) {
	p.uuid = uuid
}

func (p Post) AuthorUUID() u.UUID {
	return p.authorUUID
}

func (p *Post) SetAuthorUUID(uuid u.UUID) {
	p.authorUUID = uuid
}

func (p Post) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Post) SetCreatedAt(t time.Time) error {
	if t.After(time.Now()) {
		return ErrFutureCreatedAt
	}
	p.createdAt = t
	return nil
}

func (p Post) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Post) SetTitle(title string) {
	p.title = title
}

func (p Post) Title() string {
	return p.title
}

func (p *Post) SetContent(content string) {
	p.content = content
}

func (p Post) Content() string {
	return p.content
}

func (p Post) Likes() int {
	return p.likes
}

func (p *Post) SetLikes(likes int) error {
	if likes < 0 {
		return ErrNegativeLikes
	}
	p.likes = likes
	return nil
}

func (p *Post) IncLikes() {
	p.likes++
}

func (p *Post) SetDraft(draft bool) {
	p.draft = draft
}

func (p Post) IsDraft() bool {
	return p.draft
}

func (p *Post) Publish() error {
	if p.draft {
		return ErrPostAlreadyPublished
	}

	p.draft = false
	return nil
}

func (p *Post) SetUpdatedAt(t time.Time) error {
	if t.After(time.Now()) {
		return ErrFutureUpdatedAt
	}
	if t.Before(p.CreatedAt()) {
		return ErrInvalidUpdatedAt
	}
	p.updatedAt = t
	return nil
}

func (p Post) DeletedAt() *time.Time {
	return p.deletedAt
}

func (p *Post) SetDeletedAt(t time.Time) error {
	if t.After(time.Now()) {
		return ErrFutureDeletedAt
	}
	if t.Before(p.CreatedAt()) || t.Before(p.UpdatedAt()) {
		return ErrInvalidDeletedAt
	}
	p.deletedAt = &t
	return nil
}
