package service

import (
	"strings"
	"testing"

	"faculty/internal/model"

	"github.com/google/uuid"
)

func TestComposeComplaintAboutUser(t *testing.T) {
	targetID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	authorID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	subject, body := composeComplaint("https://faculty.example.com", complaintTarget{
		kind:        "пользователь",
		subjectNoun: "пользователя",
		name:        "Иван Иванов",
		id:          targetID,
		path:        "/students",
	}, &model.CuUserResp{ID: authorID, FirstName: "Пётр", LastName: "Петров"},
		model.ComplaintRequest{Reason: "спам", Text: "рассылает рекламу"})

	if subject != "Жалоба на пользователя: Иван Иванов" {
		t.Errorf("subject: got %q", subject)
	}

	for _, want := range []string{
		"Тип объекта: пользователь\n",
		"Объект: Иван Иванов (id: " + targetID.String() + ")\n",
		"Ссылка: https://faculty.example.com/students/" + targetID.String() + "\n",
		"Причина: спам\n",
		"Текст жалобы:\nрассылает рекламу\n",
		"Автор жалобы: Пётр Петров (id: " + authorID.String() + ")\n",
		"Профиль автора: https://faculty.example.com/students/" + authorID.String() + "\n",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("в теле письма нет %q\nтело:\n%s", want, body)
		}
	}
}

func TestComposeComplaintAboutPostWithoutBaseURL(t *testing.T) {
	postID := uuid.MustParse("33333333-3333-3333-3333-333333333333")

	subject, body := composeComplaint("", complaintTarget{
		kind:        "объявление",
		subjectNoun: "объявление",
		name:        "Ищу соседа",
		id:          postID,
		path:        "/posts",
	}, &model.CuUserResp{ID: uuid.Nil, FirstName: "Пётр", LastName: "Петров"},
		model.ComplaintRequest{Reason: "мошенничество", Text: "просит предоплату"})

	if subject != "Жалоба на объявление: Ищу соседа" {
		t.Errorf("subject: got %q", subject)
	}
	if !strings.Contains(body, "Ссылка: /posts/"+postID.String()+"\n") {
		t.Errorf("без базового URL ожидался только путь\nтело:\n%s", body)
	}
}

func TestNewComplaintServiceTrimsBaseURL(t *testing.T) {
	s := NewComplaintService(nil, nil, "https://faculty.example.com/")
	if s.frontendBaseURL != "https://faculty.example.com" {
		t.Errorf("got %q", s.frontendBaseURL)
	}
}
