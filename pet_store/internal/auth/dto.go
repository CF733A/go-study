package auth

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type TokenClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}