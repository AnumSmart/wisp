package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func (r *UserRepository) EnsureAdminExists(ctx context.Context) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		return errors.New("ADMIN_EMAIL environment variable not set")
	}

	adminPass := os.Getenv("ADMIN_PASSWORD")
	if adminPass == "" {
		return errors.New("ADMIN_PASSWORD environment variable not set")
	}

	// Проверяем наличие администратора
	exists, err := r.CheckIfInBaseByEmail(ctx, adminEmail)
	if err != nil {
		return fmt.Errorf("error checking admin existence: %w", err)
	}

	// если admin уже есть, то выходим
	if exists {
		log.Println("admin user already exists")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPass), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("admin password hashing error")
	}

	admin := &AdminConfig{
		AdminEmail: adminEmail,
		HashedPass: string(hashedPassword),
		Role:       "admin",
		Is_Active:  true,
	}

	err = r.AddUser(ctx, admin.AdminEmail, admin.HashedPass, admin.Role, admin.Is_Active)
	if err != nil {
		return err
	}
	log.Println("Admin user was added to DB")

	return nil
}
