package profile

import (
	"context"
	"errors"
	"simple_gin_server/internal/users"
	"strconv"
)

type ServiceInterface interface {
	CreateNewProfile(ctx context.Context, profile *Profile, email string) error
}

// Интерфейс для слоя ordersService для использования другими источниками
type ProfileService struct {
	repoProf ProfileRepoInterface
	repoUser users.UserRepoInterface
}

func NewProfileService(repoProf ProfileRepoInterface, repoUser users.UserRepoInterface) *ProfileService {
	return &ProfileService{
		repoProf: repoProf,
		repoUser: repoUser,
	}
}

func (p *ProfileService) CreateNewProfile(ctx context.Context, profile *Profile, email string) error {
	// Проверяем не отменен ли контекст
	if err := ctx.Err(); err != nil {
		return err
	}

	// Проверяем есть ли такой profile в базе

	profileInBase, err := p.repoProf.CheckPlrofileInBase(ctx, profile.NickName)
	if err != nil {
		return errors.New("[profile--service.go] - Failed to check profile in base")
	}

	user, err := p.repoUser.FindByEmail(ctx, email)

	if !profileInBase {
		err := p.repoProf.SaveProfile(ctx, profile, strconv.Itoa(user.Id))
		if err != nil {
			return errors.New("[profile--service.go] - Failed to save profile into base")
		}
	}

	return nil
}
