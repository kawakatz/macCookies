package types

import "time"

type Cookie struct {
	Host         string
	Path         string
	KeyName      string
	EncryptValue []byte
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

type CookieQuickManager struct {
	PathRaw           string `json:"Path raw"`
	HostRaw           string `json:"Host raw"`
	ExpiresRaw        string `json:"Expires raw"`
	ContentRaw        string `json:"Content raw"`
	NameRaw           string `json:"Name raw"`
	SameSiteRaw       string `json:"SameSite raw"`
	ThisDomainOnlyRaw string `json:"This domain only raw"`
	StoreRaw          string `json:"Store raw"`
	FirstPartyDomain  string `json:"First Party Domain"`
	HTTPOnlyRaw       string `json:"HTTP only raw,omitempty"`
}