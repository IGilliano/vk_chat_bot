package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"vk_chat_bot/pkg/storage"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Movie) error {
	fPath := filepath.Join(s.basePath, strconv.Itoa(page.UserID))

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		fmt.Println(err)
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	if err = gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userId int) (page *storage.Movie, err error) {
	path := filepath.Join(s.basePath, strconv.Itoa(userId))

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) IsExist(p *storage.Movie) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, err
	}

	path := filepath.Join(s.basePath, strconv.Itoa(p.UserID), fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func (s Storage) Remove(p *storage.Movie) error {
	fileName, err := fileName(p)
	if err != nil {
		return err
	}

	path := filepath.Join(s.basePath, strconv.Itoa(p.UserID), fileName)

	if err = os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (s Storage) decodePage(filePath string) (*storage.Movie, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var p storage.Movie

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	err = f.Close()
	if err != nil {
		return nil, err
	}
	return &p, nil

}

func fileName(p *storage.Movie) (string, error) {
	return p.Hash()
}
