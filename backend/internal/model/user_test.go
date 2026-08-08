package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRedactForBlockedViewer(t *testing.T) {
	bio := "hello"
	photo := "photos/x"
	url := "https://example.com/x"
	spec := "CS"
	status := "student"
	u := &User{
		ID:         uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		PhotoS3Key: &photo,
		PhotoURL:   &url,
		FirstName:  "Ivan",
		LastName:   "Ivanov",
		Bio:        &bio,
		Speciality: &spec,
		Status:     &status,
		Role:       "user",
	}

	RedactForBlockedViewer(u)

	if u.FirstName != "Ivan" || u.LastName != "Ivanov" {
		t.Fatalf("name cleared: %s %s", u.FirstName, u.LastName)
	}
	if u.Speciality == nil || *u.Speciality != "CS" {
		t.Fatalf("speciality: %#v", u.Speciality)
	}
	if u.PhotoS3Key != nil || u.PhotoURL != nil || u.Bio != nil || u.Status != nil {
		t.Fatalf("sensitive fields not cleared: %#v", u)
	}
	if !u.BirthDate.IsZero() {
		t.Fatal("birthdate not cleared")
	}
	if !u.BlockedByThem {
		t.Fatal("BlockedByThem not set")
	}
	if u.Role != "user" {
		t.Fatalf("role: %q", u.Role)
	}

	raw, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}
	if m["blockedByThem"] != true {
		t.Fatalf("json blockedByThem: %#v", m["blockedByThem"])
	}
	if m["birthdate"] != nil {
		t.Fatalf("json birthdate should be null, got %#v", m["birthdate"])
	}
	if m["photoUrl"] != nil {
		t.Fatalf("json photoUrl: %#v", m["photoUrl"])
	}
}

func TestRedactForDeletedAccount(t *testing.T) {
	now := time.Now()
	bio := "bio"
	u := &User{
		ID:        uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		FirstName: "Ivan",
		LastName:  "Ivanov",
		Bio:       &bio,
		Role:      "user",
		DeletedAt: &now,
	}
	if !u.IsDeleted() {
		t.Fatal("expected IsDeleted")
	}
	RedactForDeletedAccount(u)
	if !u.Deleted || u.Bio != nil || u.BlockedByThem {
		t.Fatalf("unexpected redaction: %#v", u)
	}
}
