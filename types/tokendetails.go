package types

type TokenDetails struct {
	AccessToken string
	AccessUuid  string
	//UserId       uint64
	AtExp        int64
	RefreshToken string
	RefreshUuid  string
	RtExp        int64
}
