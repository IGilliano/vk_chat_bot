package files

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"vk_chat_bot/pkg/storage"
)

var (
	ErrNoMovies = errors.New("no urls")
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) SetIntent(userID, intent int) error {
	userData, err := s.GetUserData(userID)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}

		userData.UserID = userID
	}
	userData.UserIntent = intent

	return s.Save(userData)
}

func (s Storage) Save(page *storage.UserData) error {
	fPath := filepath.Join(s.basePath, strconv.Itoa(page.UserID))

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	if err = gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	defer file.Close()

	return nil
}

func (s Storage) PickRandom(userId int) (page string, err error) {
	userData, err := s.GetUserData(userId)
	if err != nil {
		return "", err
	}
	if len(userData.Urls) == 0 {
		return "", ErrNoMovies
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(userData.Urls))

	return userData.Urls[n], nil
}

func (s Storage) IsExist(d *storage.UserData, url string) bool {
	for _, v := range d.Urls {
		if v == url {
			return true
		}
	}

	return false
}

func (s Storage) GetUserData(userID int) (*storage.UserData, error) {
	page, err := s.decodePage(filepath.Join(s.basePath, strconv.Itoa(userID)))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		page = &storage.UserData{
			UserID: userID,
		}
		err = s.Save(page)
		if err != nil {
			return nil, err
		}
	}

	return page, nil
}

func (s Storage) DeleteAll(userID int) error {
	fPath := filepath.Join(s.basePath, strconv.Itoa(userID))

	return os.Remove(fPath)
}

func (s Storage) decodePage(filePath string) (*storage.UserData, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var p storage.UserData

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return &p, nil
}
