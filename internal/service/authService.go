package service

import (
	handlers "blog-service/internal/handlers/user"
	"blog-service/internal/lib"
	"blog-service/internal/lib/hasher"
	"blog-service/internal/models"
	"blog-service/internal/repository"
	"fmt"
)

type AuthServiceImpl struct {
	userRepo     handlers.UserRepository
	tokenManager lib.TokenManager
	hasher       hasher.PasswordHasher
}

func NewAuthServiceImpl(userRepo *repository.UsersRepositoryImpl, manager *lib.Manager, hasher *hasher.SHA256Hasher) *AuthServiceImpl {
	return &AuthServiceImpl{userRepo: userRepo, tokenManager: manager, hasher: hasher}
}

func (s *AuthServiceImpl) SignIn(request models.SignInDto) (string, error) {
	const op = "service.AuthService.SignIn"
	user, err := s.userRepo.GetByCredentials(request.Login, s.hasher.Hash(request.Password))
	if err != nil {
		fmt.Printf("%s: %s\n", op, err)
		return "", err
	}
	return s.tokenManager.NewJWT(user.UserId)
}
