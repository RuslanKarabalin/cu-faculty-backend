package mailer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"mime"
	"net"
	"net/smtp"
	"strings"
	"time"
)

const (
	sendTimeout     = 15 * time.Second
	implicitTLSPort = "465"
)

type Config struct {
	Host          string
	Port          string
	Username      string
	Password      string
	From          string
	To            string
	AllowInsecure bool
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	if cfg.From == "" {
		cfg.From = cfg.Username
	}
	return &Client{cfg: cfg}
}

func (c *Client) Send(ctx context.Context, subject, body string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, sendTimeout)
		defer cancel()
	}
	deadline, _ := ctx.Deadline()

	conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", net.JoinHostPort(c.cfg.Host, c.cfg.Port))
	if err != nil {
		return fmt.Errorf("dial smtp: %w", err)
	}
	defer func() { _ = conn.Close() }()

	if err := conn.SetDeadline(deadline); err != nil {
		return fmt.Errorf("set smtp deadline: %w", err)
	}

	secure := false
	if c.cfg.Port == implicitTLSPort {
		tlsConn := tls.Client(conn, &tls.Config{ServerName: c.cfg.Host})
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			return fmt.Errorf("smtp tls handshake: %w", err)
		}
		conn = tlsConn
		secure = true
	}

	client, err := smtp.NewClient(conn, c.cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp handshake: %w", err)
	}
	defer func() { _ = client.Close() }()

	if !secure {
		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(&tls.Config{ServerName: c.cfg.Host}); err != nil {
				return fmt.Errorf("smtp starttls: %w", err)
			}
			secure = true
		}
	}
	if !secure && !c.cfg.AllowInsecure {
		return errors.New("smtp server does not support starttls")
	}

	if c.cfg.Username != "" && secure {
		auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	if err := client.Mail(c.cfg.From); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := client.Rcpt(c.cfg.To); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(buildMessage(c.cfg.From, c.cfg.To, subject, body)); err != nil {
		return fmt.Errorf("write smtp data: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("close smtp data: %w", err)
	}
	return client.Quit()
}

func buildMessage(from, to, subject, body string) []byte {
	var b strings.Builder
	b.WriteString("From: " + headerValue(from) + "\r\n")
	b.WriteString("To: " + headerValue(to) + "\r\n")
	b.WriteString("Subject: " + mime.QEncoding.Encode("utf-8", headerValue(subject)) + "\r\n")
	b.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("\r\n")
	b.WriteString(strings.ReplaceAll(strings.ReplaceAll(body, "\r\n", "\n"), "\n", "\r\n"))
	return []byte(b.String())
}

func headerValue(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}
