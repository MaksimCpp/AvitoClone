package userusecase

import (
	"context"

	"github.com/MaksimCpp/AvitoClone/internal/domain/user"
	"github.com/MaksimCpp/AvitoClone/internal/infrastructure/hash"
	jwtservice "github.com/MaksimCpp/AvitoClone/internal/infrastructure/jwt"
)

type LoginInput struct {
	Email string
	Password string
}

type LoginOutput struct {
	AccessToken string `json:"access_token"`
}

type LoginUserUseCase interface {
	Execute(ctx context.Context, input LoginInput) (*LoginOutput, error)
}

type PostgreSQLLoginUserUseCase struct {
	repo user.UserRepository
	hasher *hash.BcryptHasher
	jwtService *jwtservice.JWTService
}

func NewPostgreSQLLoginUserUseCase(
	repo user.UserRepository,
	hasher *hash.BcryptHasher,
	jwtService *jwtservice.JWTService,
) *PostgreSQLLoginUserUseCase {
	return &PostgreSQLLoginUserUseCase{
		repo: repo,
		hasher: hasher,
		jwtService: jwtService,
	}
}

func (usecase *PostgreSQLLoginUserUseCase) Execute(
	ctx context.Context, input LoginInput,
) (*LoginOutput, error) {
	userEntity, err := usecase.repo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	err = usecase.hasher.Compare(userEntity.Password, input.Password)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	token, err := usecase.jwtService.Generate(userEntity.ID)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{AccessToken: token}, nil
}
