package services

import (
	"time"

	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/types"
	"groundTurn/src/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	svc  *Services
	repo *repositories.UserRepository
}

func NewAuthService(svc *Services) *AuthService {
	return &AuthService{svc: svc, repo: repositories.NewUserRepository(svc.DB)}
}

// Login verifies bcrypt password and issues a 12h JWT.
func (s *AuthService) Login(req types.LoginRequest) (*types.LoginResponse, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, appError(401, constants.CodeAuthInvalid, "用户名或密码错误")
	}
	token, err := utils.IssueToken(s.svc.Cfg.JWTSecret, user.Username, user.Role, user.TeamID, user.Name, 12*time.Hour)
	if err != nil {
		return nil, appError(500, constants.CodeInternal, constants.MsgInternal)
	}
	return &types.LoginResponse{Token: token, User: *user}, nil
}

// EnsureSeedUsers creates the four demo accounts with hashed passwords.
func (s *AuthService) EnsureSeedUsers() error {
	count, err := s.repo.Count(s.svc.DB)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.svc.DB.Transaction(func(tx *gorm.DB) error {
		for _, acc := range constants.SeedAccounts {
			hashed, err := bcrypt.GenerateFromPassword([]byte(acc.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			user := &models.User{
				Username: acc.Username,
				Password: string(hashed),
				Name:     acc.Name,
				Role:     string(acc.Role),
				TeamID:   acc.TeamID,
			}
			if err := s.repo.Create(tx, user); err != nil {
				return err
			}
		}
		return nil
	})
}
