package models

type TokensVersion struct {
	UserID        int `json:"user_id"`
	TokensVersion int `json:"tokens_version"`
}

func (TokensVersion) TableName() string {
	return "user_tokens_versions"
}
