package domain

import "errors"

// Errors that are potentially thrown during account interactions.
var (
	ErrFutureCreatedAt      = errors.New("domain: creation time cannot be in the future")
	ErrFutureUpdatedAt      = errors.New("domain: modification time cannot be in the future")
	ErrInvalidUpdatedAt     = errors.New("domain: modification time cannot be prior to account creation time")
	ErrFutureDeletedAt      = errors.New("domain: deleton time cannot be in the future")
	ErrInvalidDeletedAt     = errors.New("domain: deletion time cannot be prior to account creation or modification time")
	ErrNegativeLikes        = errors.New("domain: likes cannot be negative")
	ErrPostAlreadyPublished = errors.New("domain: post is already published")
)
