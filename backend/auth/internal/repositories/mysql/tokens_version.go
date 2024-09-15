package mysql_repo

import (
	"gitflic.ru/project/pereverzevivan/biznes-processy-laba-1/backend/models"
	"gorm.io/gorm"
)

type TokensVersionRepo struct {
	Conn *gorm.DB
}

func NewTokensVersionRepo(conn *gorm.DB) TokensVersionRepo {
	return TokensVersionRepo{
		Conn: conn,
	}
}

func (repo TokensVersionRepo) GetUserTokensVersion(user_id int) (int, error) {

	var tokens_version int = -1
	res := repo.Conn.Model(&models.TokensVersion{}).
		Select("version").
		Where("user_id = ?", user_id).
		Scan(&tokens_version)

	if res.Error != nil {
		return -1, res.Error
	}

	return tokens_version, res.Error
}

func (repo TokensVersionRepo) IncreaseUserTokensVersion(user_id int) error {
	res := repo.Conn.Exec(`
		INSERT INTO `+models.TokensVersion{}.TableName()+` (user_id, version) 
			VALUES (?, 1)
		ON DUPLICATE KEY UPDATE version = version+1;`,
		user_id)

	return res.Error
}
