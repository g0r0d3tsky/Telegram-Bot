package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	filepath2 "path/filepath"
	"proj/lib/e"
	"proj/storage"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {

	filePath := filepath2.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, defaultPerm); err != nil {
		return e.Wrap("can't save: ", err)
	}
	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath2.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath2.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, e.Wrap("no saved files", err)
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath2.Join(path, file.Name()))
}
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page
	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("can't decode page", err)
	}
	return &p, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath2.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("can't delete file %s", path), err)
	}

	return nil
}

func (s Storage) IsExsist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file is exist", err)
	}

	path := filepath2.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't chech if file %s exists", path)
		return false, e.Wrap(msg, err)
	}
	return true, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
