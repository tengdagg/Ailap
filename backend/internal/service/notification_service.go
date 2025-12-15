package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"ailap-backend/internal/model"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// SendAlert sends a notification to the specified channel
func (s *NotificationService) SendAlert(channel *model.NotificationChannel, title, content string) error {
	var cfg map[string]string
	if err := json.Unmarshal([]byte(channel.Config), &cfg); err != nil {
		return fmt.Errorf("invalid channel config: %v", err)
	}

	if channel.Type == "webhook" {
		url := cfg["url"]
		if url == "" {
			return fmt.Errorf("webhook url is empty")
		}

		// Webhook Logic
		var payload interface{}

		// 1. Feishu / Lark (Interactive Card for Markdown)
		if strings.Contains(url, "feishu.cn") || strings.Contains(url, "larksuite.com") {
			payload = map[string]interface{}{
				"msg_type": "interactive",
				"card": map[string]interface{}{
					"header": map[string]interface{}{
						"title": map[string]string{
							"tag":     "plain_text",
							"content": title,
						},
						"template": "blue", // Header color
					},
					"elements": []interface{}{
						map[string]interface{}{
							"tag": "div",
							"text": map[string]string{
								"tag":     "lark_md",
								"content": fmt.Sprintf("%s\n\n**Time**: %s", content, time.Now().Format(time.RFC3339)),
							},
						},
						map[string]interface{}{
							"tag": "hr",
						},
						map[string]interface{}{
							"tag": "note",
							"elements": []interface{}{
								map[string]string{
									"tag":     "plain_text",
									"content": "AILAP Intelligent Monitor",
								},
							},
						},
					},
				},
			}

		} else if strings.Contains(url, "dingtalk.com") {
			// 2. DingTalk (Markdown)
			payload = map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]string{
					"title": title,
					"text":  fmt.Sprintf("### %s\n\n%s\n\n> Time: %s", title, content, time.Now().Format(time.RFC3339)),
				},
				"at": map[string]interface{}{
					"isAtAll": false,
				},
			}
		} else if strings.Contains(url, "weixin.qq.com") {
			// 3. Enterprise WeChat (Markdown)
			// QyWeixin markdown does not support 'title' field, only 'content'
			payload = map[string]interface{}{
				"msgtype": "markdown",
				"markdown": map[string]string{
					"content": fmt.Sprintf("# %s\n\n%s\n\n> Time: %s", title, content, time.Now().Format(time.RFC3339)),
				},
			}
		} else {
			// 3. Generic (Slack-compatible or simple JSON)
			// Try Slack format which many support
			payload = map[string]interface{}{
				"text": fmt.Sprintf("*%s*\n\n%s\n\nTime: %s", title, content, time.Now().Format(time.RFC3339)),
			}
		}

		body, _ := json.Marshal(payload)
		resp, err := http.Post(url, "application/json", bytes.NewReader(body))
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Read response to check for logical errors even if status is 200 (Common in DingTalk/Feishu/WeCom)
		respBody, _ := io.ReadAll(resp.Body)

		if resp.StatusCode >= 400 {
			return fmt.Errorf("webhook status %d: %s", resp.StatusCode, string(respBody))
		}

		// Check for DingTalk/WeCom error codes
		var respData map[string]interface{}
		if err := json.Unmarshal(respBody, &respData); err == nil {
			// DingTalk & WeCom use 'errcode'
			if code, ok := respData["errcode"].(float64); ok && code != 0 {
				return fmt.Errorf("webhook error: %s", string(respBody))
			}
			// Feishu uses 'code'
			if code, ok := respData["code"].(float64); ok && code != 0 {
				return fmt.Errorf("webhook error: %s", string(respBody))
			}
		}

		return nil
	}

	if channel.Type == "email" {
		host := cfg["smtp_host"]
		port := cfg["smtp_port"]
		user := cfg["username"]
		pass := cfg["password"]
		to := cfg["to"] // single recipient for simplicity, or comma separated

		if host == "" || port == "" || to == "" {
			return fmt.Errorf("incomplete email config")
		}

		addr := fmt.Sprintf("%s:%s", host, port)

		// Setup TLS config (insecure skip verify to handle user's "broken signature" error)
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}

		var client *smtp.Client
		var err error

		if port == "465" {
			// Implicit SSL/TLS for port 465
			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				return fmt.Errorf("tls dial failed: %w", err)
			}
			client, err = smtp.NewClient(conn, host)
			if err != nil {
				return fmt.Errorf("smtp new client failed: %w", err)
			}
		} else {
			// StartTLS for 587/25
			client, err = smtp.Dial(addr)
			if err != nil {
				return fmt.Errorf("smtp dial failed: %w", err)
			}
			// Try StartTLS
			if ok, _ := client.Extension("STARTTLS"); ok {
				if err = client.StartTLS(tlsConfig); err != nil {
					client.Close()
					return fmt.Errorf("starttls failed: %w", err)
				}
			}
		}
		defer client.Close()

		// Auth: Try LOGIN mechanism first as it fixes "Unrecognized authentication type" for many legacy/Exchange servers
		// Fallback to PLAIN if LOGIN is not supported or fails?
		// Actually, let's define LoginAuth helper.
		auth := LoginAuth(user, pass)

		// If server supports AUTH
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				// Retry with PLAIN if LOGIN fails? Or just return error.
				// The user error suggested PLAIN was rejected. LOGIN is likely the fix.
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}

		// Send
		if err = client.Mail(user); err != nil {
			return err
		}
		for _, recipient := range strings.Split(to, ",") {
			if strings.TrimSpace(recipient) != "" {
				if err = client.Rcpt(strings.TrimSpace(recipient)); err != nil {
					return err
				}
			}
		}

		w, err := client.Data()
		if err != nil {
			return err
		}
		msg := []byte(fmt.Sprintf("To: %s\r\n"+
			"Subject: [AILAP Alert] %s\r\n"+
			"\r\n"+
			"%s\r\n", to, title, content))

		if _, err = w.Write(msg); err != nil {
			return err
		}
		if err = w.Close(); err != nil {
			return err
		}
		return client.Quit()
	}

	return fmt.Errorf("unsupported channel type: %s", channel.Type)
}

// LoginAuth is a custom implementation for LOGIN auth mechanism
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown from server: %s", string(fromServer))
		}
	}
	return nil, nil
}
