package servicehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scheduler/internal/localdb"
)

type UserState struct {
	name     string
	date     string
	repeatly string
}

func (u *UserState) Update(key string, value string) {
	if key == "name" {
		u.name = value
	} else if key == "date" {
		u.date = value
	} else if key == "repeatly" {
		u.repeatly = value
	}
}

type MessageProcessor struct {
	localdb       *localdb.LocalDb
	userIdToState map[string]*UserState
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

func constructMessageResponse(replyToken string, message string) replyBody {
	msg := map[string]any{
		"type": "text",
		"text": message,
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
					"type":    "box",
					"layout":  "vertical",
					"spacing": "sm",
					"contents": []any{
						map[string]any{"type": "text", "text": "Name", "weight": "bold"},
						map[string]any{
							"type":    "box",
							"layout":  "horizontal",
							"spacing": "sm",
							"contents": []any{
								map[string]any{
									"type":   "button",
									"height": "sm",
									"action": map[string]any{"type": "postback", "label": "Water", "data": "job:name=Water"},
								},
								map[string]any{
									"type":   "button",
									"height": "sm",
									"action": map[string]any{"type": "postback", "label": "Pay rent", "data": "job:name=PayRent"},
								},
							},
						},
						map[string]any{
							"type":   "button",
							"height": "sm",
							"action": map[string]any{"type": "postback", "label": "Other (type)", "data": "job:name=OTHER"},
						},
					},
				},

				map[string]any{
					"type":    "box",
					"layout":  "vertical",
					"spacing": "sm",
					"contents": []any{
						map[string]any{"type": "text", "text": "Date", "weight": "bold"},
						map[string]any{
							"type":   "button",
							"height": "sm",
							"action": map[string]any{
								"type":  "datetimepicker",
								"label": "Pick date",
								"data":  "job:pick=date",
								"mode":  "date",
							},
						},
					},
				},

				map[string]any{
					"type":  "button",
					"style": "primary",
					"action": map[string]any{
						"type":  "postback",
						"label": "Submit",
						"data":  "job:submit=1",
					},
				},
			},
		},
	}
}

func NewMessageProcessor(localdb *localdb.LocalDb) *MessageProcessor {
	return &MessageProcessor{
		localdb:       localdb,
		userIdToState: map[string]*UserState{},
	}
}

func (m *MessageProcessor) getUserState(userId string) *UserState {
	userState, ok := m.userIdToState[userId]
	if !ok {
		newUserState := UserState{}
		m.userIdToState[userId] = &newUserState
		return &newUserState
	}
	return userState
}

func (m *MessageProcessor) isActionDone(userId string) bool {
	_, ok := m.userIdToState[userId]
	return !ok
}

func (m *MessageProcessor) Response(event Events) {
	replyToken := event.ReplyToken
	userId := event.Source.UserId

	if event.Type == "postback" {
		m.responsePostback(userId, replyToken, event.Postback)
		return
	}

	if event.Message.Text == "REGISTER" {
		m.messageHandler_register(userId, replyToken)
	} else if event.Message.Text == "NEW_JOB" {
		payload := constructMessageResponse(replyToken, "new")
		m.replyWtihMessage(payload)
	} else if event.Message.Text == "SHOW_ALL_JOB" {
		payload := constructMessageResponse(replyToken, "Your job")
		m.replyWtihMessage(payload)
	} else {
		payload := constructMessageResponse(replyToken, "not match")
		m.replyWtihMessage(payload)
	}
}

func (m *MessageProcessor) responsePostback(userId string, replyToken string, postback Postback) {
	userState := m.getUserState(userId)
	var msg string
	if postback.Data == "job:submit=1" {
		msg = fmt.Sprintf("Create Job: name: %s, date: %s", userState.name, userState.date)
		delete(m.userIdToState, userId)
	} else {
		for k, v := range postback.Params {
			userState.Update(k, v)
		}
		msg = fmt.Sprintf("Job: name: %s, date: %s", userState.name, userState.date)
	}
	payload := constructMessageResponse(replyToken, msg)
	m.replyWtihMessage(payload)
}

func (m *MessageProcessor) messageHandler_register(userId string, replyToken string) {
	if !m.isActionDone(userId) {
		userState := m.getUserState(userId)
		payload := constructMessageResponse(replyToken, fmt.Sprintf("Select in previous flex box, your information is name: %s, date: %v", userState.name, userState.date))
		m.replyWtihMessage(payload)
		return
	}
	userState := m.getUserState(userId)

	fmt.Println(userState.name)
	flexContent := BuildCreateJobFlex()
	payload := ConstructResponse(replyToken, flexContent, "Create Job")

	m.replyWtihMessage(payload)
}

func (m *MessageProcessor) replyWtihMessage(payload replyBody) error {

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
