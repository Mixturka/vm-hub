package config

type SessionOptions struct {
	MaxAge          int
	SessionName     string
	SessionDomain   string
	SessionSecure   bool
	SessionHttpOnly bool
	SessionFolder   string
	SessionSecret   string
	CookiesSecret   string
}
