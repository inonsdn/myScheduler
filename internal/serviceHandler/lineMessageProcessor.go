package servicehandler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"scheduler/internal/localdb"
)

type MessageProcessor struct {
	localdb *localdb.LocalDb
}

type replyBody struct {
	ReplyToken string        `json:"replyToken"`
	Messages   []interface{} `json:"messages"`
}

type MessageContent struct {
	Type string      `json:"type"`
	Body MessageBody `json:"body"`
}
type MessageBody struct {
	Type     string           `json:"type"`
	Layout   string           `json:"layout"`
	Spacing  string           `json:"spacing"`
	Contents []MessageContent `json:"contents"`
}

func ConstructResponse(replyToken string, flexContents any, altText string) replyBody {
	msg := map[string]any{
		"type":     "flex",
		"altText":  altText,
		"contents": flexContents, // bubble JSON object
	}

	return replyBody{
		ReplyToken: replyToken,
		Messages:   []interface{}{msg},
	}
}

func BuildCreateJobFlex() map[string]any {
	return map[string]any{
		"type": "bubble",
		"body": map[string]any{
			"type":    "box",
			"layout":  "vertical",
			"spacing": "md",
			"contents": []any{
				map[string]any{"type": "text", "text": "Create Job", "weight": "bold", "size": "xl"},
				map[string]any{
					"type":  "button",
					"style": "primary",
					"action": map[string]any{
						"type":  "postback",
						"label": "Submit",
						"data":  "job:submit=1",
					},
				},
				// add your date/time/repeat buttons here...
			},
		},
	}
}

func NewMessageProcessor(localdb *localdb.LocalDb) *MessageProcessor {
	return &MessageProcessor{
		localdb: localdb,
	}
}

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
	// payload := map[string]any{
	// 	"replyToken": replyToken,
	// 	"messages": []map[string]any{
	// 		{"type": "text", "text": message},
	// 	},
	// }

	flexContent := BuildCreateJobFlex()
	payload := ConstructResponse(replyToken, flexContent, message)

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
