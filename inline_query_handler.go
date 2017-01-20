package openpoll

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/davecgh/go-spew/spew"
	"golang.org/x/net/context"
	"gopkg.in/kyokomi/emoji.v1"
	"github.com/google/uuid"
)

type ResultArticleContent struct {
	title    string
	variants []string
}

var resultArticlesContent = []ResultArticleContent{
	{
		title:    "Start :thumbsup: / :thumbsdown: poll",
		variants: []string{":thumbsup:", ":thumbsdown:"},
	},
	{
		title:    "Start :thumbsup: / :neutral_face: / :thumbsdown: poll",
		variants: []string{":thumbsup:", ":neutral_face:", ":thumbsdown:"},
	},
	{
		title:    "Start :ok_hand: poll",
		variants: []string{":ok_hand:"},
	},
}

func HandleInlineQuery(ctx context.Context, query *telegram.InlineQuery) (err error) {
	api := telebot.GetAPI(ctx)

	resultQueryArticles := []telegram.InlineQueryResult{}

	for _, resultArticleContent := range resultArticlesContent {
		pollUuid, err := uuid.NewRandom()
		pollId := pollUuid.String()
		if err != nil {
			return err
		}

		buttons := []telegram.InlineKeyboardButton{}

		for i, replyVariant := range resultArticleContent.variants {
			button := telegram.InlineKeyboardButton{
				CallbackData: fmt.Sprintf("p:%v:%v", pollId, i),
				Text:         emoji.Sprint(replyVariant),
			}
			buttons = append(buttons, button)
		}

		replyMarkup := telegram.InlineKeyboardMarkup{
			InlineKeyboard: [][]telegram.InlineKeyboardButton{
				buttons,
			},
		}

		resultArticle := telegram.InlineQueryResultArticle{
			BaseInlineQueryResult: telegram.BaseInlineQueryResult{
				Type: "article",
				ID:   pollId,
				InputMessageContent: telegram.InputTextMessageContent{
					MessageText: fmt.Sprintf("*%v*", query.Query),
					ParseMode:   "Markdown",
				},
				ReplyMarkup: replyMarkup,
			},
			Title: emoji.Sprintf(resultArticleContent.title),
		}

		// TODO: Register poll in redis

		resultQueryArticles = append(resultQueryArticles, resultArticle)
	}

	cfg := telegram.AnswerInlineQueryCfg{
		InlineQueryID: query.ID,
		CacheTime:     0,
		IsPersonal:    true,
		Results:       resultQueryArticles,
	}

	_, err = api.AnswerInlineQuery(ctx, cfg)
	if err != nil {
		logrus.Info("Failed to send inline query answer:", spew.Sdump(err))
	}
	return err
}
