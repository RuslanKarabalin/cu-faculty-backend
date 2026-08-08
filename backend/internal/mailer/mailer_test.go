package mailer

import (
	"slices"
	"strings"
	"testing"
)

func TestBuildMessageHeaders(t *testing.T) {
	msg := string(buildMessage("from@example.com", "to@example.com", "Жалоба на пользователя: Иван Иванов", "строка 1\nстрока 2"))

	head, body, ok := strings.Cut(msg, "\r\n\r\n")
	if !ok {
		t.Fatal("в сообщении нет пустой строки между заголовками и телом")
	}

	lines := strings.Split(head, "\r\n")
	for _, want := range []string{
		"From: from@example.com",
		"To: to@example.com",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
	} {
		if !slices.Contains(lines, want) {
			t.Errorf("нет заголовка %q", want)
		}
	}

	if !slices.ContainsFunc(lines, func(line string) bool { return strings.HasPrefix(line, "Date: ") }) {
		t.Error("нет заголовка Date")
	}

	if strings.Contains(head, "Иван") {
		t.Error("тема должна быть закодирована по RFC2047, а не отправлена как есть")
	}
	if !strings.Contains(head, "Subject: =?utf-8?q?") {
		t.Errorf("тема не закодирована: %q", head)
	}
	if body != "строка 1\r\nстрока 2" {
		t.Errorf("тело: got %q", body)
	}
}

func TestBuildMessageStripsHeaderInjection(t *testing.T) {
	msg := string(buildMessage("from@example.com", "to@example.com", "тема\r\nBcc: attacker@example.com", "тело"))

	head, _, _ := strings.Cut(msg, "\r\n\r\n")
	lines := strings.Split(head, "\r\n")
	if len(lines) != 7 {
		t.Fatalf("ожидалось 7 строк заголовков, got %d: %q", len(lines), head)
	}
	for _, line := range lines {
		if strings.HasPrefix(line, "Bcc") {
			t.Errorf("перевод строки в теме позволил подделать заголовок: %q", head)
		}
	}
}
