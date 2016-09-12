package main

import (
	"fmt"
	"github.com/dveselov/go-mystem"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"strings"
)

var (
	Bot                  *tgbotapi.BotAPI
	TelegramBotToken     = os.Getenv("TELEGRAM_BOT_TOKEN")
	PrettyGrammemesTable = map[int]string{
		mystem.Substantive: "сущ.",
		mystem.Verb:        "глаг.",
	}
)

type Result struct {
	form      string
	grammemes []int
}

func (result *Result) String() string {
	return fmt.Sprintf("%s - %v", result.form, result.PrettyGrammemes())
}

func (result *Result) Markdown() string {
	return fmt.Sprintf("***%s*** - `%s`", result.form, result.PrettyGrammemes())
}

func (result *Result) PrettyGrammemes() string {
	var (
		results []string
	)
	for _, grammeme := range result.grammemes {
		pretty, ok := PrettyGrammemesTable[grammeme]
		if !ok {
			pretty = "неизв."
		}
		results = append(results, pretty)
	}
	return strings.Join(results, ", ")
}

func handleInlineQuery(query *tgbotapi.InlineQuery) {
	log.Printf("[%s] %s", query.From, query.Query)
	var (
		articles []interface{}
	)
	analyses := mystem.NewAnalyses(query.Query)
	for i := 0; i < analyses.Count(); i++ {
		lemma := analyses.GetLemma(i)
		id := fmt.Sprintf("%s:%d", query.ID, i)
		result := Result{lemma.Form(), lemma.StemGram()}
		article := tgbotapi.NewInlineQueryResultArticleMarkdown(id, result.String(), result.Markdown())
		articles = append(articles, article)
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       articles,
	}
	if _, err := Bot.AnswerInlineQuery(inlineConf); err != nil {
		log.Println(err)
	}
	defer analyses.Close()
}

func main() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(TelegramBotToken)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Authorized on account %s.", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery == nil {
			continue
		}
		if update.InlineQuery.Query != "" {
			go handleInlineQuery(update.InlineQuery)
		}
	}
}
