package services

import (
	"groundTurn/src/constants"
	"groundTurn/src/models"
	"groundTurn/src/repositories"
	"groundTurn/src/utils"
	"time"

	"gorm.io/gorm"
)

// AuthService issues local JWTs from the seeded app_user table.
type AuthService struct {
	Users  *repositories.UserRepository
	secret string
}

func NewAuthService(db *gorm.DB, secret string) *AuthService {
	return &AuthService{Users: repositories.NewUserRepository(db), secret: secret}
}

func (a *AuthService) Login(username string) (string, *models.User, error) {
	user, err := a.Users.GetByUsername(username)
	if err != nil {
		return "", nil, utils.NewAPIError(401, constants.CodeAuthRequired, "用户不存在："+username, nil)
	}
	token, err := utils.IssueToken(a.secret, user.Username, user.Role, user.TeamID, 12*time.Hour)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}
