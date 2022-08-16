package models

import "time"

type TokenDetails struct {
	AccessToken  string
	AccessUuid   string
	UserId       uint64 `json:"id"`
	Username     string `json:"user_name"`
	AtExp        time.Time
	RefreshToken string
	RefreshUuid  string
	RtExp        int64
}
