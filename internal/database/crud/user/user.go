package user

import (
	"context"
	"errors"

	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/database/schema"
	"gorm.io/gorm"
)

var (
	ErrUsernameEmpty = errors.New("empty user name")
	ErrEmailEmpty    = errors.New("empty email address")
)

func CreateUser(db *gorm.DB, ctx context.Context, username, email string) error {
	user := schema.User{
		Username:    username,
		Email:       email,
		Displayname: username,
	}
	if err := validateUserData(&user); err != nil {
		return err
	}
	return gorm.G[schema.User](db).Create(ctx, &user)
}

func GetUserByID(db *gorm.DB, ctx context.Context, id uint) (*schema.User, error) {
	user, err := gorm.G[schema.User](db).Where("id = ?", id).First(ctx)
	return &user, err
}

func GetUserByEmail(db *gorm.DB, ctx context.Context, email string) (*schema.User, error) {
	user, err := gorm.G[schema.User](db).Where("email = ?", email).First(ctx)
	return &user, err
}

func GetUsers(db *gorm.DB, ctx context.Context) ([]schema.User, error) {
	return gorm.G[schema.User](db).Find(ctx)
}

func UpdateUser(db *gorm.DB, ctx context.Context, user *schema.User) error {
	if err := validateUserData(user); err != nil {
		return err
	}
	_, err := gorm.G[schema.User](db).Where("id = ?", user.ID).Updates(ctx, *user)
	return err
}

func DeleteUser(db *gorm.DB, ctx context.Context, id uint) error {
	_, err := gorm.G[schema.User](db).Where("id = ?", id).Delete(ctx)
	return err
}

func validateUserData(user *schema.User) error {
	if user.Username == "" {
		return ErrUsernameEmpty
	}
	if user.Email == "" {
		return ErrEmailEmpty
	}
	return nil
}
