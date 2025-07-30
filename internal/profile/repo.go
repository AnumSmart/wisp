package profile

import (
	"context"
	"fmt"
	"log"
	"simple_gin_server/pkg/db"

	"github.com/google/uuid"
)

// Интерфейс для слоя ordersRepository для использования другими источниками
type ProfileRepoInterface interface {
	CheckPlrofileInBase(ctx context.Context, nickName string) (bool, error)
	SaveProfile(ctx context.Context, profile *Profile, userId string) error
}

type ProfileRepository struct {
	Database db.PgRepoInterface
}

// Конструктор репозитория
func NewProfileRepository(dataBase db.PgRepoInterface) *ProfileRepository {
	return &ProfileRepository{
		Database: dataBase,
	}
}

// Метод репоизтория profile для проверки наличия профиля по уникальному nickName
func (p *ProfileRepository) CheckPlrofileInBase(ctx context.Context, nickName string) (bool, error) {

	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return false, err
	}

	var exists bool
	// Используем COUNT и проверяем, есть ли хотя бы одна запись
	query := `SELECT EXISTS(SELECT 1 FROM profiles WHERE nick_name = $1)`
	err := p.Database.GetPool().QueryRow(ctx, query, nickName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check profile existence: %w", err)
	}

	return exists, nil
}

// Метод репоизтория profile для записи профиля в базу
func (p *ProfileRepository) SaveProfile(ctx context.Context, profile *Profile, userId string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	query := `INSERT INTO profiles (id, user_id, name, nick_name, gender, age_group, city, profession, smoking, goal, hobbies, social_link, rating) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := p.Database.GetPool().Exec(ctx, query,
		uuid.New().String(),
		profile.UserID,
		profile.Name,
		profile.NickName,
		profile.Gender,
		profile.AgeGroup,
		profile.City,
		profile.Profession,
		profile.Smoking,
		profile.Goal,
		profile.Hobbies,
		profile.SocialLink,
		10,
	)

	if err != nil {
		log.Printf("[profile--repo.go]--Error during saving profile to DB:%v", err.Error())
		return err
	}
	return nil
}
