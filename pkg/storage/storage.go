package storage

import (
	"errors"
)

type Storage interface {
	Save(p *UserData) error
	PickRandom(userId int) (string, error)
	IsExist(p *UserData, url string) bool
	SetIntent(userID, intent int) error
	GetUserData(userID int) (*UserData, error)
	DeleteAll(userID int) error
}

var ErrNoSavedPages = errors.New("no saved pages")

type UserData struct {
	Urls       []string
	UserID     int
	UserIntent int
}

type User struct {
	MoviesCount int
	Intent      int
}
