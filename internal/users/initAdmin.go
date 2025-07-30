package users

import (
	"context"
	"errors"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func InitAdmin(ctx context.Context, repo UserRepoInterface) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")

	isInBase, _ := repo.CheckIfInBaseByEmail(ctx, adminEmail)

	if !isInBase {
		adminPass := os.Getenv("ADMIN_PASSWORD")
		// Хеширование пароля
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("ошибка при хешировании пароля админа")
		}

		err = repo.AddUser(ctx, adminEmail, string(hashedPassword), "admin", true)
		if err != nil {
			return err
		}
	}
	return nil
}
