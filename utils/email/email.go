package email

import (
	"github.com/yuanJewel/go-core/logger"
	"io"
	"time"
)

type Body struct {
	Id      string
	From    string
	Cc      string
	To      string
	Subject string
	Date    time.Time
	Body    string
}

func GetAllEmails(option Opt, startTime, endTime time.Time) ([]Body, error) {
	c, err := newPop3Client(option).NewConn()
	if err != nil {
		return nil, err
	}
	defer func(c *Conn) {
		_ = c.Quit()
	}(c)

	var body []Body
	count, _, _ := c.Stat()
	for id := count - 1; id > 0; id-- {
		m, err := c.Retr(id)
		if err != nil {
			return nil, err
		}

		subject, err := parseSubject(m.Header.Get("subject"))
		if err != nil {
			logger.Log.Warnf("parse subject error: %v", err)
			continue
		}

		parsedTime, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", m.Header.Get("date"))
		if err != nil {
			logger.Log.Warnf("parse date error: %v", err)
			continue
		}

		if parsedTime.After(endTime) {
			continue
		}

		if parsedTime.Before(startTime) {
			break
		}

		data, err := io.ReadAll(m.Body)
		if err != nil {
			logger.Log.Warnf("read body error: %v", err)
			continue
		}

		body = append(body, Body{
			Id:      m.Header.Get("message-id"),
			From:    m.Header.Get("from"),
			Cc:      m.Header.Get("cc"),
			To:      m.Header.Get("to"),
			Date:    parsedTime,
			Subject: subject,
			Body:    string(data),
		})
	}
	return body, nil
}
