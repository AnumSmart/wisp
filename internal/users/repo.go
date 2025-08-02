package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"simple_gin_server/pkg/db"

	"github.com/jackc/pgx"
)

// Интерфейс для слоя userRepository для использования другими источниками
type UserRepoInterface interface {
	AddUser(ctx context.Context, email, hashedPass string, role string, is_active bool) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	GetEmailLIst(ctx context.Context) ([]string, error)
	CheckIfInBaseByEmail(ctx context.Context, email string) (bool, error)
	AddRefreshToken(ctx context.Context, email, refreshToken string) error
	ClearRefreshToken(ctx context.Context, claimsEmail string) error
	EnsureAdminExists(ctx context.Context) error
}

type UserRepository struct {
	Database db.PgRepoInterface
}

// Конструктор репозитория
func NewUserRepository(dataBase db.PgRepoInterface) *UserRepository {
	return &UserRepository{
		Database: dataBase,
	}
}

// Сохранение пользователя (с хешированным паролем)
func (r *UserRepository) AddUser(ctx context.Context, email, hashedPass string, role string, is_active bool) error {

	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	// создаем экземпляр нового юзера для сохранения в БД
	newUser := User{
		Email:    email,
		HashPass: hashedPass,
		Role:     role,
		IsActive: is_active,
	}
	log.Println(newUser.Email)
	query := `INSERT INTO users (email, hashed_pass, user_role, is_active) VALUES ($1, $2, $3, $4) ON CONFLICT (email) DO NOTHING`
	res, err := r.Database.GetPool().Exec(ctx, query, newUser.Email, newUser.HashPass, newUser.Role, newUser.IsActive)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("email already registered")
	}

	return nil
}

// Проверка наличия пользователя в хранилище
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const query = `
		SELECT id, email, hashed_pass, COALESCE(refresh_token, '')
		FROM users 
		WHERE email = $1
		LIMIT 1
	`
	var user User
	err := r.Database.GetPool().QueryRow(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.HashPass,
		&user.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%v: %v", ErrUserNotExists, err)
		}
		log.Printf("function [FindByEmail], failed to query user by email: %v", err)
		return nil, errors.New(ErrUserNotExists)
	}

	return &user, nil
}

// Получение слайса email зарегестрированных пользователей
func (r *UserRepository) GetEmailLIst(ctx context.Context) ([]string, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	// создаём строку запроса в БД
	const query = `SELECT email FROM users`

	// получаем результат запроса
	rows, err := r.Database.GetPool().Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query email list: %w", err)
	}
	// закрываем соединение
	defer rows.Close()

	var emails []string

	// начинаем вычитывать по 1 элементу из каждого ряда
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, fmt.Errorf("failed to scan email: %w", err)
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %w", err)
	}

	return emails, nil
}

// Проверка того, что пользователь с указанным email уже находится в БД
func (r *UserRepository) CheckIfInBaseByEmail(ctx context.Context, email string) (bool, error) {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return false, err
	}

	const query = `
        SELECT EXISTS(
            SELECT 1 
            FROM users 
            WHERE email = $1
        )
    `

	var exists bool
	err := r.Database.GetPool().QueryRow(ctx, query, email).Scan(&exists)

	if err != nil {
		log.Printf("ошибка при запросе в базу на поиск по Email: %v", err)
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}

	return exists, nil
}

// добавляем поле refreshToken в базу по email (нужно держать refreshToken в БД)
func (r *UserRepository) AddRefreshToken(ctx context.Context, email, refreshToken string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	const query = `
		UPDATE users 
		SET refresh_token = $1 
		WHERE email = $2;
	`
	_, err := r.Database.GetPool().Exec(ctx, query, refreshToken, email)
	if err != nil {
		log.Printf("[repo.go]---[AddRefreshToken()]---Err: %v", err)
		return err
	}

	return nil
}

// удаляет refresh токен из базы данных по заданному Email из claims
func (r *UserRepository) ClearRefreshToken(ctx context.Context, claimsEmail string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	const query = `
		UPDATE users 
		SET refresh_token = NULL 
		WHERE email = $1;
	`
	_, err := r.Database.GetPool().Exec(ctx, query, claimsEmail)
	if err != nil {
		return err
	}

	return nil
}
