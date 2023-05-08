package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Storage interface {
	Save(p *Movie) error
	PickRandom(userId int) (*Movie, error)
	Remove(p *Movie) error
	IsExist(p *Movie) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Movie struct {
	Url    string
	UserID int
}

func (p Movie) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.Url); err != nil {
		fmt.Println(err)
		return "", err
	}

	if _, err := io.WriteString(h, strconv.Itoa(p.UserID)); err != nil {
		fmt.Println(err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
