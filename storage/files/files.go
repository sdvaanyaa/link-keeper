package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"saveBot/lib/errwrap"
	"saveBot/storage"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0174

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = errwrap.WrapIfErr("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err = gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errwrap.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return errwrap.Wrap("can't remove page", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)
	if err = os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove page file %s:", path)
		return errwrap.Wrap(msg, err)
	}
	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, errwrap.Wrap("can't check if file exist", err)
	}
	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file exists: %s", path)
		return false, errwrap.Wrap(msg, err)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (page *storage.Page, err error) {
	defer func() { err = errwrap.WrapIfErr("can't decode page", err) }()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() { _ = file.Close() }()

	var p storage.Page

	if err = gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}
