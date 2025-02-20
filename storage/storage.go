package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	"saveBot/lib/errwrap"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(*Page) error
	IsExists(p *Page) (bool, error)
}

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
