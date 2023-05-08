package vk

import (
	"encoding/json"
	"fmt"
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
	case 0:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Старт",
			}, Color: "positive"},
			},
		}
		return makeKeyboard(buttons)

	case 1:
		buttons := [][]Button{
			{Button{Action: Action{
				Type:  "text",
				Label: "Посоветуй фильм!",
			}, Color: "positive"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Оценить фильм",
			}, Color: "primary"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "Моя коллекция",
			}, Color: "primary"},
			},
			{Button{Action: Action{
				Type:  "text",
				Label: "SOS!",
			}, Color: "primary"},
			},
		}
		return makeKeyboard(buttons)
	case 2:
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
	case 3:
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
	case 4:
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
