package internal

import (
	"github.com/adlio/trello"
)

type Trello struct {
	client *trello.Client
}

func NewTrello(key string, token string) Trello {
	client := trello.NewClient(key, token)
	client.Logger = logger

	return Trello{
		client: client,
	}
}

func (t Trello) GetMember() *trello.Member {
	member, err := t.client.GetMember("me", trello.Defaults())

	if err != nil {
		logger.WithError(err).WithField("member", "me").Fatalf("failed to get member")
	}

	return member
}

func (t Trello) GetMemberBoards(member *trello.Member) []*trello.Board {
	boards, err := member.GetBoards(trello.Arguments{
		"filter":              "all",
		"fields":              "all",
		"lists":               "all",
		"organization":        "true",
		"organization_fields": "all",
	})

	if err != nil {
		logger.WithError(err).WithField("memberID", member.ID).Fatalf("failed to get user boards")
	}

	return boards
}

func (t Trello) GetListCards(list *trello.List) []*trello.Card {
	args := trello.Arguments{
		"attachments":       "true",
		"attachment_fields": "all",
		"customFieldItems":  "true",
		"stickers":          "true",
		"members":           "true",
		"member_fields":     "all",
		"checkItemStates":   "true",
		"checklists":        "all",
		"limit":             "1000",
		"sort":              "-id",
		"filter":            "all",
		"fields":            "all",
		"pluginData":        "true",
	}

	cards, err := list.GetCards(args)

	if err != nil {
		logger.WithError(err).WithField("listID", list.ID).Fatalf("failed to get list cards")
	}

	return cards
}
