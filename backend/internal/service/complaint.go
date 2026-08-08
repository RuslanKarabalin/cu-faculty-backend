package service

import (
	"context"
	"strings"

	"faculty/internal/model"

	"github.com/google/uuid"
)

type complaintMailer interface {
	Send(ctx context.Context, subject, body string) error
}

type complaintRepository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetPostByID(ctx context.Context, id uuid.UUID) (*model.Post, error)
}

type ComplaintService struct {
	repo            complaintRepository
	mailer          complaintMailer
	frontendBaseURL string
}

func NewComplaintService(repo complaintRepository, mailer complaintMailer, frontendBaseURL string) *ComplaintService {
	return &ComplaintService{repo: repo, mailer: mailer, frontendBaseURL: strings.TrimRight(frontendBaseURL, "/")}
}

type complaintTarget struct {
	kind        string
	subjectNoun string
	name        string
	id          uuid.UUID
	path        string
}

func (s *ComplaintService) ComplainAboutUser(ctx context.Context, complainant *model.CuUserResp, targetID uuid.UUID, req model.ComplaintRequest) error {
	user, err := s.repo.GetUserByID(ctx, targetID)
	if err != nil {
		return err
	}

	subject, body := composeComplaint(s.frontendBaseURL, complaintTarget{
		kind:        "пользователь",
		subjectNoun: "пользователя",
		name:        personName(user.FirstName, user.LastName),
		id:          user.ID,
		path:        "/students",
	}, complainant, req)
	return s.mailer.Send(ctx, subject, body)
}

func (s *ComplaintService) ComplainAboutPost(ctx context.Context, complainant *model.CuUserResp, targetID uuid.UUID, req model.ComplaintRequest) error {
	post, err := s.repo.GetPostByID(ctx, targetID)
	if err != nil {
		return err
	}

	subject, body := composeComplaint(s.frontendBaseURL, complaintTarget{
		kind:        "объявление",
		subjectNoun: "объявление",
		name:        post.Title,
		id:          post.ID,
		path:        "/posts",
	}, complainant, req)
	return s.mailer.Send(ctx, subject, body)
}

func composeComplaint(baseURL string, target complaintTarget, complainant *model.CuUserResp, req model.ComplaintRequest) (string, string) {
	subject := "Жалоба на " + target.subjectNoun + ": " + target.name

	var b strings.Builder
	b.WriteString("Тип объекта: " + target.kind + "\n")
	b.WriteString("Объект: " + target.name + " (id: " + target.id.String() + ")\n")
	b.WriteString("Ссылка: " + complaintLink(baseURL, target.path, target.id) + "\n")
	b.WriteString("Причина: " + req.Reason + "\n")
	b.WriteString("Текст жалобы:\n")
	b.WriteString(req.Text + "\n")
	b.WriteString("---\n")
	b.WriteString("Автор жалобы: " + personName(complainant.FirstName, complainant.LastName) + " (id: " + complainant.ID.String() + ")\n")
	b.WriteString("Профиль автора: " + complaintLink(baseURL, "/students", complainant.ID) + "\n")
	return subject, b.String()
}

func complaintLink(baseURL, path string, id uuid.UUID) string {
	return baseURL + path + "/" + id.String()
}

func personName(firstName, lastName string) string {
	return strings.TrimSpace(firstName + " " + lastName)
}
