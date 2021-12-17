package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/adlio/trello"
)

// BackupAllBoards backup all boards to directory
func BackupAllBoards(t Trello, backupTo string) {
	boards := t.GetMemberBoards(t.GetMember())

	for _, board := range boards {
		boardDir := resolveWhereToSaveTheBoard(backupTo, board, true)

		for _, list := range board.Lists {
			listDir := resolveWhereToSaveList(boardDir, list, true)

			cards := t.GetListCards(list)
			for _, card := range cards {
				cardDir := resolveWhereToSaveCard(listDir, card, true)
				downloadCardAttachments(t.client.Key, t.client.Token, cardDir, card.Attachments)
			}

			saveCardsToJSONFile(listDir, cards)
		}

		saveListsToJSONFile(boardDir, board.Lists)
	}

	saveBoardsToJSONFile(backupTo, boards)
}

func downloadCardAttachments(key string, token string, cardDir string, attachments []*trello.Attachment) {
	for _, attachment := range attachments {
		url := attachment.URL

		logger.WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).Debugln("downloading attachment")

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			logger.WithError(err).WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).Fatalln("unable to create request")
		}

		query := req.URL.Query()
		query.Add("key", key)
		query.Add("token", token)
		req.URL.RawQuery = query.Encode()

		req.Header.Add("Authorization", fmt.Sprintf("OAuth oauth_consumer_key=\"%s\", oauth_token=\"%s\"", key, token))

		client := &http.Client{}
		resp, err := client.Do(req.WithContext(context.Background()))
		if err != nil {
			logger.WithError(err).WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).Fatalln("unable to download resource")
		}
		defer func() {
			if err = resp.Body.Close(); err != nil {
				logger.WithError(err).WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).Debugln("unable to create request")
			}
		}()

		filename := fmt.Sprintf("%s%c%s", cardDir, os.PathSeparator, Sanitize(attachment.Name))

		// Create the file
		out, err := os.Create(filename)
		if err != nil {
			logger.WithError(err).WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).WithField("filename", filename).Debugln("unable to create file on the destination path")
		}
		defer func() {
			if err = out.Close(); err != nil {
				logger.WithError(err).WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).WithField("filename", filename).Debugln("unable to close the file")
			}
		}()

		_, err = io.Copy(out, resp.Body)

		if err != nil {
			logger.WithError(err).WithField("cardDir", cardDir).WithField("attachmentURL", attachment.URL).WithField("filename", filename).Fatalln("unable to copy the content of the downloaded file to the local file")
		}
	}
}

func saveCardsToJSONFile(saveTo string, cards []*trello.Card) {
	marshalIndent, err2 := json.MarshalIndent(cards, "", "  ")

	if err2 != nil {
		logger.WithError(err2).WithField("cards", func() []string {
			var result []string
			for _, card := range cards {
				result = append(result, card.Name)
			}
			return result
		}).Fatalf("failed to marshal cards")
	}

	SaveFile(fmt.Sprintf("%s%ccards.json", saveTo, os.PathSeparator), marshalIndent)
}

func saveListsToJSONFile(saveTo string, lists []*trello.List) {
	marshalIndent, err2 := json.MarshalIndent(lists, "", "  ")

	if err2 != nil {
		logger.WithError(err2).WithField("lists", func() []string {
			var result []string
			for _, list := range lists {
				result = append(result, list.Name)
			}
			return result
		}).Fatalf("failed to marshal lists")
	}

	SaveFile(fmt.Sprintf("%s%clists.json", saveTo, os.PathSeparator), marshalIndent)
}

func saveBoardsToJSONFile(saveTo string, boards []*trello.Board) {
	marshalIndent, err := json.MarshalIndent(boards, "", "  ")

	if err != nil {
		logger.WithError(err).WithField("boards", func() []string {
			var result []string
			for _, board := range boards {
				result = append(result, board.Name)
			}
			return result
		}).Fatalf("failed to marshal boards")
	}

	SaveFile(fmt.Sprintf("%s%cboards.json", saveTo, os.PathSeparator), marshalIndent)
}

func resolveWhereToSaveCard(listDir string, card *trello.Card, createDir bool) string {
	resourceName := Sanitize(card.Name)
	if card.Closed {
		resourceName = addNamePrefixToClosedResource(resourceName)
	}

	dir := filepath.Join(listDir, resourceName)
	if createDir {
		CreateDirectoryRecursively(dir)
	}

	return dir
}

func resolveWhereToSaveList(boardDir string, list *trello.List, createDir bool) string {
	resourceName := Sanitize(list.Name)
	if list.Closed {
		resourceName = addNamePrefixToClosedResource(resourceName)
	}

	dir := filepath.Join(boardDir, resourceName)
	if createDir {
		CreateDirectoryRecursively(dir)
	}

	return dir
}

func resolveWhereToSaveTheBoard(backupTo string, board *trello.Board, createDir bool) string {
	resourceName := Sanitize(board.Name)
	if board.Closed {
		resourceName = addNamePrefixToClosedResource(resourceName)
	}

	dir := filepath.Join(backupTo, resourceName)
	if createDir {
		CreateDirectoryRecursively(dir)
	}

	return dir
}

func addNamePrefixToClosedResource(resourceName string) string {
	return fmt.Sprintf("[closed] %s", resourceName)
}
