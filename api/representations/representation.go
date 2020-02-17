package representations

import (
	"time"

	"github.com/freerware/negotiator/representation"
)

// Representation is the core representation.
type Representation struct {
	representation.Base

	lastModified *time.Time
	etag         *string
}

func (r Representation) LastModified() *time.Time { return r.lastModified }

func (r *Representation) SetLastModified(t time.Time) { r.lastModified = &t }

func (r Representation) ETag() *string { return r.etag }

func (r *Representation) SetETag(tag string) { r.etag = &tag }
