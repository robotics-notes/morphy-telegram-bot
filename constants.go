package main

import "github.com/dveselov/go-mystem"

var (
	PrettyGrammemesTable = map[int]string{
		mystem.Substantive:  "сущ.",
		mystem.Verb:         "глаг.",
		mystem.Feminine:     "жен.",
		mystem.Animated:     "оживл.",
		mystem.FirstName:    "имя",
		mystem.Imperfect:    "несов.",
		mystem.Intransitive: "нп.",
	}
)
