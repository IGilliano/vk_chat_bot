package vk

import (
	"encoding/json"
	"fmt"
)

type MenuType int

const (
	MenuTypeEmpty = iota
	MenuTypeMain
	MenuTypeRecommendation
	MenuTypeSos
	MenuTypeCollection
	MenuTypeCollectionAdd
	MenuTypeStart
)

type KeyBoard struct {
	OneTime bool       `json:"one_time"`
	Buttons [][]Button `json:"buttons"`
}

type Button struct {
	Action Action `json:"action"`
	Color  string `json:"color"`
}

type Action struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label"`
}

func makeButtons(num int) string {
	switch num {
	case MenuTypeStart:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Старт",
			}, Color: "positive"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Очистить данные",
			}, Color: "negative"},
			},
		}
		return makeKeyboard(buttons)

	case MenuTypeMain:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Посоветуй фильм",
			}, Color: "positive"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Моя коллекция",
			}, Color: "primary"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Помощь",
			}, Color: "primary"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Стоп",
			}, Color: "negative"},
			},
		}
		return makeKeyboard(buttons)
	case MenuTypeRecommendation:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Спасибо!",
			}, Color: "positive"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Хочу другой",
			}, Color: "negative"},
			},
		}
		return makeKeyboard(buttons)
	case MenuTypeSos:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Как добавлять новые фильмы?",
			}, Color: "primary"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "В главное меню",
			}, Color: "primary"},
			},
		}
		return makeKeyboard(buttons)
	case MenuTypeCollection:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Добавить фильм",
			}, Color: "positive"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Удалить фильм",
			}, Color: "negative"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "В главное меню",
			}, Color: "primary"},
			},
		}
		return makeKeyboard(buttons)
	case MenuTypeCollectionAdd:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Добавить фильм",
			}, Color: "positive"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "В главное меню",
			}, Color: "primary"},
			},
		}
		return makeKeyboard(buttons)
	default:
		return ""
	}
}

func makeKeyboard(buttons [][]Button) string {
	kb := KeyBoard{
		OneTime: true,
		Buttons: buttons,
	}

	keyboard, err := json.Marshal(kb)
	if err != nil {
		fmt.Printf("cant parse keyboard %s", err)
	}
	return string(keyboard)
}
