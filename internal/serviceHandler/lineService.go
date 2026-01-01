package servicehandler

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"scheduler/internal/config"
)

var (
	channelSecret = os.Getenv("CHANNEL_SECRET")
	accessToken   = os.Getenv("CHANNEL_ACCESS_TOKEN")
)

type LineService struct {
	webhookUrl    string
	port          int
	channelSecret string
	accessToken   string
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

// Example for payload from line
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

func verifySignature(body []byte, got string) bool {
	mac := hmac.New(sha256.New, []byte(channelSecret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	// LINE uses base64; compare timing-safe
	return hmac.Equal([]byte(expected), []byte(got))
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, _ := io.ReadAll(r.Body)
	sig := r.Header.Get("x-line-signature")
	if sig == "" || channelSecret == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !verifySignature(body, sig) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var webhook Webhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Reply quickly; LINE expects 200 fast
	w.WriteHeader(http.StatusOK)

	for _, event := range webhook.events {
		replyMessage(event.ReplyToken, "test")
	}
}

func NewLineService(opts *config.Options) *LineService {
	lineOpts := opts.GetLineOptions()
	return &LineService{
		webhookUrl:    lineOpts.GetWebhookUrl(),
		port:          lineOpts.GetPort(),
		channelSecret: lineOpts.GetChannelSecret(),
		accessToken:   lineOpts.GetAccessToken(),
	}
}

func (l *LineService) Run() {
	http.HandleFunc(l.webhookUrl, webhookHandler)
	runningPort := fmt.Sprintf("0.0.0.0:%d", l.port)
	fmt.Println("Run serve at ", runningPort)
	http.ListenAndServe(runningPort, nil)
}

func (l *LineService) RegisterRoute() {
}

func (l *LineService) OnShutdown() {
}
