package utils

import "github.com/gofiber/fiber/v3/log"

// Логирует ошибку, если она есть, а также возвращает nil
func LogErrorIfNotEmpty(err error) error {
	if err.Error() != "" {
		log.Error(err)
	}
	return nil
}
