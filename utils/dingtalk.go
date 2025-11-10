package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type dingTalkMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At                 dingTalkAt         `json:"at"`
	DingTalkActionCard DingTalkActionCard `json:"actionCard"`
}

type dingTalkAt struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type DingTalkActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
}

func sendDingTalk(webhook string, message dingTalkMessage) error {
	// 将消息结构体转为 JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// 如果请求成功，返回结果
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("failed to send message, status code: %d, send dingTalk Body is %s, return is %s",
		resp.StatusCode, string(messageJSON), string(respBytes))
}

func SendDingTalkMessage(webhook string, content string, at []string, isAtAll bool) error {
	message := dingTalkMessage{
		MsgType: "text",
		At: dingTalkAt{
			AtMobiles: at,
			IsAtAll:   isAtAll,
		},
	}
	if len(at) != 0 {
		content += fmt.Sprintf("\n通知人: @%s", strings.Join(at, " @"))
	}
	message.Text.Content = content

	return sendDingTalk(webhook, message)
}

func SendDingTalkActionCard(webhook string, actionCard DingTalkActionCard, at []string, isAtAll bool) error {
	if actionCard.Title == "" {
		actionCard.BtnOrientation = "0"
	} else {
		actionCard.Text = fmt.Sprintf("## %s \n\n %s", actionCard.Title, actionCard.Text)
	}
	if actionCard.BtnOrientation == "" {
		actionCard.BtnOrientation = "0"
	}
	if actionCard.SingleTitle == "" {
		actionCard.SingleTitle = "点击查看"
	}

	actionCard.SingleURL = fmt.Sprintf("dingtalk://dingtalkclient/page/link?pc_slide=false&url=%s",
		url.QueryEscape(actionCard.SingleURL))
	if at != nil && len(at) != 0 {
		actionCard.Text += fmt.Sprintf("\n\n通知人: @%s", strings.Join(at, " @"))
	}
	message := dingTalkMessage{
		MsgType:            "actionCard",
		DingTalkActionCard: actionCard,
		At: dingTalkAt{
			AtMobiles: at,
			IsAtAll:   isAtAll,
		},
	}

	return sendDingTalk(webhook, message)
}
