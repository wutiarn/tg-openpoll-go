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

var expirationTime = 60*60*24*10

func HandleInlineQuery(ctx context.Context, query *telegram.InlineQuery) (err error) {
	api := telebot.GetAPI(ctx)
	redis, err := GetRedisConnection()
	if err != nil {
		return err
	}

	resultQueryArticles := []telegram.InlineQueryResult{}

	for _, resultArticleContent := range resultArticlesContent {
		pollUuid, err := uuid.NewRandom()
		pollId := pollUuid.String()
		if err != nil {
			return err
		}

		pollTopic := fmt.Sprintf("*%v*", query.Query)

		topicKey := fmt.Sprintf("openpoll_%s_topic", pollId)
		variantsKey := fmt.Sprintf("openpoll_%s_variant", pollId)

		redis.Do("SET", topicKey, pollTopic)
		redis.Do("EXPIRE", topicKey, expirationTime)

		buttons := []telegram.InlineKeyboardButton{}

		for i, replyVariant := range resultArticleContent.variants {
			button := telegram.InlineKeyboardButton{
				CallbackData: fmt.Sprintf("p:%v:%v", pollId, i),
				Text:         emoji.Sprint(replyVariant),
			}
			buttons = append(buttons, button)

			redis.Do("RPUSH", variantsKey, replyVariant)
		}
		redis.Do("EXPIRE", variantsKey)

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
					MessageText: pollTopic,
					ParseMode:   "Markdown",
				},
				ReplyMarkup: replyMarkup,
			},
			Title: emoji.Sprintf(resultArticleContent.title),
		}

		logrus.Info("Created poll %v for query %v", pollId, query.ID)

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
