package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"saveBot/lib/errwrap"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
	// Created time.Time
}

func (p Page) Hash() (hash string, err error) {
	defer func() { err = errwrap.WrapIfErr("can't calculate hash", err) }()
	h := sha1.New()
	if _, err = io.WriteString(h, p.URL); err != nil {
		return "", err
	}
	if _, err = io.WriteString(h, p.UserName); err != nil {
		return "", err
	}
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return hash, nil
}
