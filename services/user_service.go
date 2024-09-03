package services

import (
	"pharmacy-backend/models"

	"gorm.io/gorm"
)

type UserService struct {
    db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
    var users []models.User
    if err := s.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}

func (s *UserService) CreateUser(user *models.User) error {
    if err := s.db.Create(user).Error; err != nil {
        return err
    }
    return nil
}

func (s *UserService) GetUserByEmail(email string) (models.User, error) {
    var user models.User
    if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
        return user, err
    }

    return user, nil
}

