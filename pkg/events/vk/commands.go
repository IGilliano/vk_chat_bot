package vk

import (
	"errors"
	"log"
	"net/url"
	"strings"
	"vk_chat_bot/pkg/storage"
)

func (p *Processor) doCmd(text string, userId int) error {
	text = strings.TrimSpace(text)

	log.Printf("Got new message %s from user:%d", text, userId)

	//добавить/удалить?
	if isUrl(text) {
		return p.saveMovie(userId, text)
	}

	switch text {
	case "Привет", "Старт":
		return p.hello(userId)
	case "Посоветуй фильм!", "Хочу другой":
		return p.sendRandom(userId)
	case "Моя коллекция":
		return p.vk.SendMessage(userId, msgCollection, makeButtons(4))
	case "Добавить фильм":
		return p.vk.SendMessage(userId, msgSaveMovie, makeButtons(4))
	case "Удалить фильм":
		return p.vk.SendMessage(userId, msgSaveMovie, makeButtons(4))
	case "SOS!":
		return p.help(userId)
	case "Как добавлять новые фильмы?":
		return p.vk.SendMessage(userId, msgHelp1, makeButtons(1))
	case "Спасибо!":
		return p.vk.SendMessage(userId, msgTY, makeButtons(1))
	case "В главное меню":
		return p.vk.SendMessage(userId, msgMainMenu, makeButtons(1))
	case "Стоп":
		return p.vk.SendMessage(userId, msgStop, makeButtons(0))
	default:
		return p.vk.SendMessage(userId, msgUnknownCommand, makeButtons(1))
	}
}

func (p *Processor) saveMovie(userId int, movieUrl string) error {

	page := &storage.Movie{
		Url:    movieUrl,
		UserID: userId,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.vk.SendMessage(userId, msgAlreadyExists, makeButtons(1))
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err = p.vk.SendMessage(userId, msgSaved, makeButtons(1)); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(userId int) error {
	page, err := p.storage.PickRandom(userId)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	keyboard := makeButtons(2)

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.vk.SendMessage(userId, msgNoSavedMovies, makeButtons(1))
	}

	if err = p.vk.SendMessage(userId, page.Url, keyboard); err != nil {
		return err
	}

	return nil
}

func (p *Processor) deleteMovie(userId int, movieUrl string) error {
	page := &storage.Movie{
		Url:    movieUrl,
		UserID: userId,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.vk.SendMessage(userId, msgAlreadyNotExists, makeButtons(4))
	}

	err = p.storage.Remove(page)
	if err != nil {
		return err
	}

	if err = p.vk.SendMessage(userId, msgDelete, makeButtons(4)); err != nil {
		return err
	}

	return nil
}

func (p *Processor) help(userId int) error {
	return p.vk.SendMessage(userId, msgHelp, makeButtons(3))

}

func (p *Processor) hello(userId int) error {
	return p.vk.SendMessage(userId, msgHello, makeButtons(1))
}

func isUrl(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
