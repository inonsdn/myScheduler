package servicehandler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type MessageProcessor struct{}

func (m *MessageProcessor) Response(event Events) {
	replyToken := event.ReplyToken
	if event.Message.Text == "REGISTER" {
		m.replyWtihMessage(replyToken, "registered")
	} else if event.Message.Text == "NEW_JOB" {
		m.replyWtihMessage(replyToken, "registered")
	} else if event.Message.Text == "SHOW_ALL_JOB" {
		m.replyWtihMessage(replyToken, "Here is registered job:")
	} else {
		m.replyWtihMessage(replyToken, "Not match")
	}
}

func (m *MessageProcessor) replyWtihMessage(replyToken string, message string) error {
	payload := map[string]any{
		"replyToken": replyToken,
		"messages": []map[string]any{
			{"type": "text", "text": message},
		},
	}

	b, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.line.me/v2/bot/message/reply", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		body, _ := io.ReadAll(resp.Body)
		return &httpError{Status: resp.Status, Body: string(body)}
	}
	return nil
}
