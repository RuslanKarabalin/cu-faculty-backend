package model

type ComplaintRequest struct {
	Reason string `json:"reason"`
	Text   string `json:"text"`
}
