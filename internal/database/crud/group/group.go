package group

import (
	"context"
	"errors"

	userCrud "github.com/Bokume2/FirstHackathon2026Summer_back/internal/database/crud/user"
	"github.com/Bokume2/FirstHackathon2026Summer_back/internal/database/schema"
	"gorm.io/gorm"
)

var (
	ErrNonexistentUser = errors.New("nonexistent user")
)

func CreateGroup(db *gorm.DB, ctx context.Context, name string, user_ids []string) (*schema.Group, error) {
	member := make([]*schema.User, len(user_ids))
	for _, id := range user_ids {
		user, err := userCrud.GetUserByID(db, ctx, id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNonexistentUser
		} else if err != nil {
			return nil, err
		}
		member = append(member, user)
	}
	group := schema.Group{
		Name:   name,
		Member: member,
	}
	err := gorm.G[schema.Group](db).Omit("Member").Create(ctx, &group)
	return &group, err
}

func GetGroupByID(db *gorm.DB, ctx context.Context, id uint) (*schema.Group, error) {
	group, err := gorm.G[schema.Group](db).Where("id = ?", id).First(ctx)
	return &group, err
}

func UpdateGroup(db *gorm.DB, ctx context.Context, group *schema.Group) error {
	_, err := gorm.G[schema.Group](db).Where("id = ?").Updates(ctx, *group)
	return err
}

func AddUsersToGroup(db *gorm.DB, ctx context.Context, id uint, user_ids []string) (*schema.Group, error) {
	group, err := GetGroupByID(db, ctx, id)
	if err != nil {
		return nil, err
	}
	newMember := make([]*schema.User, len(user_ids))
	for _, id := range user_ids {
		user, err := userCrud.GetUserByID(db, ctx, id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNonexistentUser
		} else if err != nil {
			return nil, err
		}
		newMember = append(newMember, user)
	}
	group.Member = append(group.Member, newMember...)
	err = UpdateGroup(db, ctx, group)
	return group, err
}
