package email

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

type SMTPOpt struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	Log      *logrus.Entry
}

func (c *SMTPOpt) SendMail(to, cc []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer func() {
		if closeErr := conn.Close(); closeErr != nil {
			c.Log.Warnf("关闭连接时出错: %v", closeErr)
		}
	}()

	reader := bufio.NewReader(conn)
	smtp := smtpConn{conn, reader, c.Log}

	if err = smtp.checkConnect(); err != nil {
		return err
	}

	if err = smtp.auth(c.Username, c.Password); err != nil {
		return err
	}

	if err = smtp.sendFrom(c.From); err != nil {
		return err
	}

	if err = smtp.sendTo(to); err != nil {
		return err
	}

	if err := smtp.sendCc(cc); err != nil {
		return err
	}

	var headers []string
	headers = append(headers, fmt.Sprintf("From: %s", c.From))
	if len(to) > 0 {
		headers = append(headers, fmt.Sprintf("To: %s", strings.Join(to, ", ")))
	}
	if len(cc) > 0 {
		headers = append(headers, fmt.Sprintf("Cc: %s", strings.Join(cc, ", ")))
	}
	headers = append(headers, fmt.Sprintf("Subject: %s", subject))
	headers = append(headers, "MIME-Version: 1.0", "Content-Type: text/html; charset=utf-8", "")

	bodies := []string{"<html><body>"}
	bodies = append(bodies, "<div style=\"font-family: Arial, Helvetica, sans-serif; font-size: 16px; line-height: 1.4; color: #333333;\">")
	bodies = append(bodies, "<br>", body, "</div></body></html>")
	message := fmt.Sprintf("%s\r\n%s\r\n.\r\n", strings.Join(headers, "\r\n"), strings.Join(bodies, "\r\n"))
	c.Log.Debugf("发送邮件内容: %s", message)

	if err = smtp.sendData(message); err != nil {
		return err
	}

	return smtp.quit()
}

type smtpConn struct {
	conn   net.Conn
	reader *bufio.Reader
	Log    *logrus.Entry
}

func (c *smtpConn) checkConnect() error {
	welcome, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("读取欢迎信息失败: %v", err)
	}
	c.Log.Debugf("服务器响应: %s", welcome)

	if !strings.HasPrefix(welcome, "220") {
		return fmt.Errorf("服务器返回错误: %s", welcome)
	}

	return nil
}

func (c *smtpConn) auth(username, password string) error {
	if username == "" && password == "" {
		c.Log.Debug("未提供用户名和密码，将使用匿名登录")
		return nil
	}

	authMethods, err := c.getAuthMethods()
	if err != nil {
		return err
	}
	if len(authMethods) == 0 {
		return fmt.Errorf("服务器不支持任何认证方法")
	}
	c.Log.Debugf("服务器支持的认证方法: %v", authMethods)

	if len(authMethods) > 0 && username != "" && password != "" {
		authenticated := false

		if contains(authMethods, "PLAIN") {
			err := c.authPlain(username, password)
			if err == nil {
				authenticated = true
				c.Log.Debug("PLAIN认证成功")
			} else {
				c.Log.Debugf("PLAIN认证失败: %v", err)
			}
		}

		if !authenticated && contains(authMethods, "LOGIN") {
			err := c.authLogin(username, password)
			if err == nil {
				authenticated = true
				c.Log.Debug("LOGIN认证成功")
			} else {
				c.Log.Debugf("LOGIN认证失败: %v", err)
			}
		}

		if !authenticated {
			return fmt.Errorf("所有系统录入的认证方法都失败了")
		}
	}
	return nil
}

func (c *smtpConn) getAuthMethods() ([]string, error) {
	authMethods := []string{}
	if _, err := fmt.Fprintf(c.conn, "EHLO localhost\r\n"); err != nil {
		return authMethods, fmt.Errorf("发送EHLO命令失败: %v", err)
	}
	for {
		response, err := c.reader.ReadString('\n')
		if err != nil {
			return authMethods, fmt.Errorf("读取EHLO响应失败: %v", err)
		}
		c.Log.Debugf("EHLO响应: %s", response)

		if strings.HasPrefix(response, "250-") {
			if strings.Contains(response, "AUTH") {
				authParts := strings.Split(response, "AUTH")
				if len(authParts) > 1 {
					methods := strings.Split(strings.TrimSpace(authParts[1]), " ")
					authMethods = append(authMethods, methods...)
				}
			}
		} else if strings.HasPrefix(response, "250 ") {
			break
		}
	}
	return authMethods, nil
}

func (c *smtpConn) authPlain(username, password string) error {
	authString := fmt.Sprintf("\000%s\000%s", username, password)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(authString))

	if _, err := fmt.Fprintf(c.conn, "AUTH PLAIN %s\r\n", authEncoded); err != nil {
		return fmt.Errorf("发送AUTH PLAIN命令失败: %v", err)
	}
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "235") {
		return fmt.Errorf("PLAIN认证失败: %s", response)
	}

	return nil
}

func (c *smtpConn) authLogin(username, password string) error {
	if _, err := fmt.Fprintf(c.conn, "AUTH LOGIN\r\n"); err != nil {
		return fmt.Errorf("发送AUTH LOGIN命令失败: %v", err)
	}
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "334") {
		return fmt.Errorf("AUTH LOGIN命令被拒绝: %s", response)
	}

	usernameEncoded := base64.StdEncoding.EncodeToString([]byte(username))
	if _, err := fmt.Fprintf(c.conn, "%s\r\n", usernameEncoded); err != nil {
		return fmt.Errorf("发送用户名失败: %v", err)
	}
	response, err = c.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "334") {
		return fmt.Errorf("用户名被拒绝: %s", response)
	}

	passwordEncoded := base64.StdEncoding.EncodeToString([]byte(password))
	if _, err := fmt.Fprintf(c.conn, "%s\r\n", passwordEncoded); err != nil {
		return fmt.Errorf("发送密码失败: %v", err)
	}
	response, err = c.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if !strings.HasPrefix(response, "235") {
		return fmt.Errorf("密码被拒绝: %s", response)
	}

	return nil
}

func (c *smtpConn) sendFrom(from string) error {
	if _, err := fmt.Fprintf(c.conn, "MAIL FROM:<%s>\r\n", from); err != nil {
		return fmt.Errorf("发送MAIL FROM命令失败: %v", err)
	}
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("MAIL FROM命令失败: %v", err)
	}
	c.Log.Debugf("MAIL FROM响应: %s", response)

	if !strings.HasPrefix(response, "250") {
		return fmt.Errorf("MAIL FROM命令被拒绝: %s", response)
	}
	return nil
}

func (c *smtpConn) sendTo(to []string) error {
	for _, recipient := range to {
		if _, err := fmt.Fprintf(c.conn, "RCPT TO:<%s>\r\n", recipient); err != nil {
			return fmt.Errorf("发送RCPT TO命令失败: %v", err)
		}
		response, err := c.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("RCPT TO命令失败: %v", err)
		}
		c.Log.Debugf("RCPT TO响应: %s", response)

		if !strings.HasPrefix(response, "250") {
			return fmt.Errorf("RCPT TO命令被拒绝: %s", response)
		}
	}
	return nil
}

func (c *smtpConn) sendCc(cc []string) error {
	for _, recipient := range cc {
		if _, err := fmt.Fprintf(c.conn, "RCPT TO:<%s>\r\n", recipient); err != nil {
			return fmt.Errorf("发送RCPT TO命令失败: %v", err)
		}
		response, err := c.reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("RCPT TO命令失败: %v", err)
		}
		c.Log.Debugf("RCPT TO响应(Cc): %s", response)

		if !strings.HasPrefix(response, "250") {
			return fmt.Errorf("RCPT TO命令被拒绝: %s", response)
		}
	}
	return nil
}

func (c *smtpConn) sendData(message string) error {
	if _, err := fmt.Fprintf(c.conn, "DATA\r\n"); err != nil {
		return fmt.Errorf("发送DATA命令失败: %v", err)
	}
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("DATA命令失败: %v", err)
	}
	c.Log.Debugf("DATA响应: %s", response)

	if !strings.HasPrefix(response, "354") {
		return fmt.Errorf("DATA命令被拒绝: %s", response)
	}

	if _, err := fmt.Fprintf(c.conn, message); err != nil {
		return fmt.Errorf("发送邮件内容失败: %v", err)
	}
	response, err = c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("发送邮件内容失败: %v", err)
	}
	c.Log.Debugf("邮件发送响应: %s", response)

	if !strings.HasPrefix(response, "250") {
		return fmt.Errorf("邮件发送被拒绝: %s", response)
	}

	return nil
}

func (c *smtpConn) quit() error {
	if _, err := fmt.Fprintf(c.conn, "QUIT\r\n"); err != nil {
		return fmt.Errorf("发送QUIT命令失败: %v", err)
	}
	response, err := c.reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("QUIT命令失败: %v", err)
	}
	c.Log.Debugf("QUIT响应: %s", response)
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
