package servicehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
// channelSecret = os.Getenv("CHANNEL_SECRET")
// accessToken   = os.Getenv("CHANNEL_ACCESS_TOKEN")
)

type LineService struct {
}

type Webhook struct {
	events []Events
}

type Events struct {
	ReplyToken string  `json:"replyToken"`
	Type       string  `json:"type"`
	Mode       string  `json:"mode"`
	Timestamp  int     `json:"timestamp"`
	Source     Source  `json:"source"`
	EventId    string  `json:"eventId"`
	Message    Message `json:"message"`
}

type Source struct {
	Type    string `json:"type"`
	GroupId string `json:"groupId"`
	UserId  string `json:"userId"`
}

type Message struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	QuoteToken      string `json:"quoteToken"`
	MarkAsReadToken string `json:"markAsReadToken"`
	Text            string `json:"text"`
}

// // When a user sends a text message containing mention and an emoji in a group chat
// {
//   "destination": "xxxxxxxxxx",
//   "events": [
//     {
//       "replyToken": "nHuyWiB7yP5Zw52FIkcQobQuGDXCTA",
//       "type": "message",
//       "mode": "active",
//       "timestamp": 1462629479859,
//       "source": {
//         "type": "group",
//         "groupId": "Ca56f94637c...",
//         "userId": "U4af4980629..."
//       },
//       "webhookEventId": "01FZ74A0TDDPYRVKNK77XKC3ZR",
//       "deliveryContext": {
//         "isRedelivery": false
//       },
//       "message": {
//         "id": "444573844083572737",
//         "type": "text",
//         "quoteToken": "q3Plxr4AgKd...",
//         "markAsReadToken": "30yhdy232...",
//         "text": "@All @example Good Morning!! (love)",
//         "emojis": [
//           {
//             "index": 29,
//             "length": 6,
//             "productId": "5ac1bfd5040ab15980c9b435",
//             "emojiId": "001"
//           }
//         ],
//         "mention": {
//           "mentionees": [
//             {
//               "index": 0,
//               "length": 4,
//               "type": "all"
//             },
//             {
//               "index": 5,
//               "length": 8,
//               "userId": "U49585cd0d5...",
//               "type": "user",
//               "isSelf": false
//             }
//           ]
//         }
//       }
//     }
//   ]
// }

type httpError struct {
	Status string
	Body   string
}

func (h httpError) Error() string {
	return fmt.Sprintf("Error with status %s msg: %s", h.Status, h.Body)
}

func replyMessage(replyToken string, message string) error {
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

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	var webhook Webhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, event := range webhook.events {
		replyMessage(event.ReplyToken, "test")
	}

}

func (l *LineService) Run() {
	http.HandleFunc("/line/webhook", webhookHandler)
}
