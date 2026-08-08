package model

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	maxBioLen             = 255
	maxSpecialityLen      = 63
	maxShortTitleLen      = 63
	maxTitleLen           = 127
	maxContentLen         = 255
	maxPlaceLen           = 63
	maxSocialLinkLen      = 127
	maxRegLinkLen         = 255
	maxCompanyNameLen     = 63
	maxPositionLen        = 63
	maxEduLevelLen        = 31
	maxSpecializionLen    = 63
	maxNoteLen            = 255
	maxComplaintReasonLen = 127
	maxComplaintTextLen   = 1000
	maxPublishDays        = 3650
	minValidYear          = 1900
)

func maxValidYear() int { return time.Now().Year() + 10 }

type ValidationError struct{ Message string }

func (e *ValidationError) Error() string { return e.Message }

type fieldErrors struct{ msgs []string }

func (f *fieldErrors) add(format string, args ...any) {
	f.msgs = append(f.msgs, fmt.Sprintf(format, args...))
}

func (f *fieldErrors) required(field, val string) {
	if strings.TrimSpace(val) == "" {
		f.add("%s is required", field)
	}
}

func (f *fieldErrors) maxLen(field, val string, n int) {
	if utf8.RuneCountInString(val) > n {
		f.add("%s must be at most %d characters", field, n)
	}
}

func (f *fieldErrors) requiredMax(field, val string, n int) {
	f.required(field, val)
	f.maxLen(field, val, n)
}

func (f *fieldErrors) optionalMax(field string, val *string, n int) {
	if val != nil {
		f.maxLen(field, *val, n)
	}
}

func (f *fieldErrors) year(field string, y int) {
	if y < minValidYear || y > maxValidYear() {
		f.add("%s must be between %d and %d", field, minValidYear, maxValidYear())
	}
}

func (f *fieldErrors) yearRange(startYear int, endYear *int) {
	f.year("startYear", startYear)
	if endYear != nil {
		f.year("endYear", *endYear)
		if *endYear < startYear {
			f.add("endYear must not be before startYear")
		}
	}
}

func (f *fieldErrors) result() error {
	if len(f.msgs) == 0 {
		return nil
	}
	return &ValidationError{Message: strings.Join(f.msgs, "; ")}
}

func (r UpdateUserRequest) Validate() error {
	var f fieldErrors
	f.optionalMax("bio", r.Bio, maxBioLen)
	f.optionalMax("speciality", r.Speciality, maxSpecialityLen)
	if r.StatusID != nil && *r.StatusID <= 0 {
		f.add("statusId must be a positive integer")
	}
	return f.result()
}

func (r CreatePostRequest) Validate() error {
	var f fieldErrors
	f.requiredMax("title", r.Title, maxShortTitleLen)
	f.requiredMax("content", r.Content, maxContentLen)
	return f.result()
}

func (r UpdatePostRequest) Validate() error {
	var f fieldErrors
	f.requiredMax("title", r.Title, maxShortTitleLen)
	f.requiredMax("content", r.Content, maxContentLen)
	return f.result()
}

func (r SocialRequest) Validate() error {
	var f fieldErrors
	f.required("social", r.Social)
	f.requiredMax("link", r.Link, maxSocialLinkLen)
	return f.result()
}

func (r WorkPlaceRequest) Validate() error {
	var f fieldErrors
	f.requiredMax("companyName", r.CompanyName, maxCompanyNameLen)
	f.required("grade", r.Grade)
	f.requiredMax("position", r.Position, maxPositionLen)
	f.yearRange(r.StartYear, r.EndYear)
	return f.result()
}

func (r EduPlaceRequest) Validate() error {
	var f fieldErrors
	if r.UniversityId <= 0 {
		f.add("universityId is required")
	}
	f.required("grade", r.Grade)
	f.optionalMax("level", r.Level, maxEduLevelLen)
	f.requiredMax("specialization", r.Specialization, maxSpecializionLen)
	f.yearRange(r.StartYear, r.EndYear)
	return f.result()
}

func (r CreateNewsRequest) Validate() error {
	var f fieldErrors
	f.requiredMax("title", r.Title, maxTitleLen)
	f.requiredMax("content", r.Content, maxContentLen)
	if r.PublishDays > maxPublishDays {
		f.add("publishDays must be at most %d", maxPublishDays)
	}
	return f.result()
}

func (r UpdateNewsRequest) Validate() error {
	var f fieldErrors
	if r.Title != nil {
		f.requiredMax("title", *r.Title, maxTitleLen)
	}
	if r.Content != nil {
		f.requiredMax("content", *r.Content, maxContentLen)
	}
	if r.PublishDays != nil && *r.PublishDays > maxPublishDays {
		f.add("publishDays must be at most %d", maxPublishDays)
	}
	return f.result()
}

func (r CreateEventRequest) Validate() error {
	var f fieldErrors
	f.requiredMax("title", r.Title, maxTitleLen)
	f.requiredMax("content", r.Content, maxContentLen)
	f.requiredMax("place", r.Place, maxPlaceLen)
	f.required("category", r.Category)
	if r.StartsAt.IsZero() {
		f.add("startsAt is required")
	}
	f.optionalMax("registrationLink", r.RegistrationLink, maxRegLinkLen)
	return f.result()
}

func (r UpdateEventRequest) Validate() error {
	var f fieldErrors
	if r.Title != nil {
		f.requiredMax("title", *r.Title, maxTitleLen)
	}
	if r.Content != nil {
		f.requiredMax("content", *r.Content, maxContentLen)
	}
	if r.Place != nil {
		f.requiredMax("place", *r.Place, maxPlaceLen)
	}
	if r.Category != nil {
		f.required("category", *r.Category)
	}
	if r.StartsAt != nil && r.StartsAt.IsZero() {
		f.add("startsAt must be a valid time")
	}
	f.optionalMax("registrationLink", r.RegistrationLink, maxRegLinkLen)
	return f.result()
}

func (r CreateContactRequest) Validate() error {
	var f fieldErrors
	if r.ContactID == (uuid.UUID{}) {
		f.add("contactId is required")
	}
	f.optionalMax("note", r.Note, maxNoteLen)
	return f.result()
}

func (r UpdateContactRequest) Validate() error {
	var f fieldErrors
	f.optionalMax("note", r.Note, maxNoteLen)
	return f.result()
}

func (r ComplaintRequest) Validate() error {
	var f fieldErrors
	f.requiredMax("reason", r.Reason, maxComplaintReasonLen)
	f.requiredMax("text", r.Text, maxComplaintTextLen)
	return f.result()
}
