package vk

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"vk_chat_bot/pkg/storage"
	"vk_chat_bot/pkg/storage/files"
)

func (p *Processor) doCmd(text string, userId int) error {
	text = strings.TrimSpace(text)

	log.Printf("Got new message %s from user:%d", text, userId)

	if u, err := url.Parse(text); err == nil && u.Host != "" {
		fmt.Println(u.Hostname())
		if u.Hostname() != "kinopoisk.ru" && u.Hostname() != "www.kinopoisk.ru" {
			return p.vk.SendMessage(userId, msgBotIsStupid, makeButtons(MenuTypeCollection))
		}

		return p.processMovie(userId, text)
	}

	var err error

	switch text {
	case "Привет", "Старт":
		err = p.hello(userId)
	case "Посоветуй фильм", "Хочу другой":
		err = p.sendRandom(userId)
	case "Моя коллекция":
		err = p.myCollection(userId)
	case "Добавить фильм":
		return p.setIntent(userId, 1)
	case "Удалить фильм":
		return p.setIntent(userId, 2)
	case "Помощь":
		err = p.help(userId)
	case "Как добавлять новые фильмы?":
		err = p.vk.SendMessage(userId, msgHelp1, makeButtons(MenuTypeMain))
	case "Спасибо!":
		err = p.vk.SendMessage(userId, msgTY, makeButtons(MenuTypeMain))
	case "В главное меню":
		err = p.vk.SendMessage(userId, msgMainMenu, makeButtons(MenuTypeMain))
	case "Стоп":
		err = p.vk.SendMessage(userId, msgStop, makeButtons(MenuTypeStart))
	case "Очистить данные":
		err = p.deleteAll(userId)
	default:
		err = p.vk.SendMessage(userId, msgUnknownCommand, makeButtons(MenuTypeMain))
	}
	if err != nil {
		return err
	}

	return p.setIntent(userId, 0)
}

func (p *Processor) myCollection(userId int) error {
	userData, err := p.storage.GetUserData(userId)
	if err != nil {
		return err
	}

	if len(userData.Urls) == 0 {
		return p.vk.SendMessage(userId, msgNoSavedMovies, makeButtons(MenuTypeCollectionAdd))
	}
	var subResult string
	for i, v := range userData.Urls {
		subResult += fmt.Sprintf(msgCollectionPartFmt, i+1, v)
	}
	result := fmt.Sprintf(msgCollectionMainFmt, subResult)

	err = p.storage.Save(userData)
	if err != nil {
		return err
	}

	return p.vk.SendMessage(userId, result, makeButtons(MenuTypeCollection))
}

func (p *Processor) setIntent(userId, intent int) error {
	if err := p.storage.SetIntent(userId, intent); err != nil {
		return err
	}

	var msg string
	switch intent {
	case 1:
		msg = msgSaveMovie
	case 2:
		msg = msgDeleteMovie
	default:
		return nil
	}
	return p.vk.SendMessage(userId, msg, makeButtons(MenuTypeEmpty))
}

func (p *Processor) processMovie(userId int, movieUrl string) error {
	userData, err := p.storage.GetUserData(userId)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	switch userData.UserIntent {
	case 0, 1:
		return p.addMovie(userData, movieUrl)
	case 2:
		return p.deleteMovie(userData, movieUrl)
	default:
		return errors.New("something unexpected happened")
	}
}

func (p *Processor) addMovie(userData *storage.UserData, url string) error {
	for _, v := range userData.Urls {
		if v == url {
			return p.vk.SendMessage(userData.UserID, msgAlreadyExists, makeButtons(MenuTypeMain))
		}
	}
	userData.Urls = append(userData.Urls, url)
	userData.UserIntent = 0

	if err := p.storage.Save(userData); err != nil {
		return err
	}

	return p.vk.SendMessage(userData.UserID, msgSaved, makeButtons(MenuTypeMain))
}

func (p *Processor) deleteMovie(userData *storage.UserData, movieUrl string) error {
	if !p.storage.IsExist(userData, movieUrl) {
		return p.vk.SendMessage(userData.UserID, msgAlreadyNotExists, makeButtons(MenuTypeMain))
	}

	newUrls, err := deleteUrl(userData.Urls, movieUrl)
	if err != nil {
		return err
	}

	userData.UserIntent = 0
	userData.Urls = newUrls

	if err = p.storage.Save(userData); err != nil {
		return err
	}

	return p.vk.SendMessage(userData.UserID, msgDelete, makeButtons(MenuTypeMain))
}

func (p *Processor) sendRandom(userId int) error {
	page, err := p.storage.PickRandom(userId)
	if err != nil {
		if !errors.Is(err, files.ErrNoMovies) {
			return err
		}

		return p.vk.SendMessage(userId, msgNoSavedMovies, makeButtons(MenuTypeCollectionAdd))
	}

	return p.vk.SendMessage(userId, page, makeButtons(MenuTypeRecommendation))
}

func deleteUrl(urls []string, toDelete string) ([]string, error) {
	newUrls := make([]string, 0, len(urls)-1)
	for _, v := range urls {
		if v == toDelete {
			continue
		}
		newUrls = append(newUrls, v)
	}

	return newUrls, nil
}

func (p *Processor) deleteAll(userId int) error {
	if err := p.storage.DeleteAll(userId); err != nil {
		return err
	}

	return p.vk.SendMessage(userId, msgDeletedAll, makeButtons(MenuTypeEmpty))
}

func (p *Processor) help(userId int) error {
	return p.vk.SendMessage(userId, msgHelp, makeButtons(MenuTypeSos))

}

func (p *Processor) hello(userId int) error {
	return p.vk.SendMessage(userId, msgHello, makeButtons(MenuTypeMain))
}
