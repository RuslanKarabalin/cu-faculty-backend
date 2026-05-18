package cuclient

import "errors"

const CookieName = "bff.cookie"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrUpstream     = errors.New("upstream error")
)
