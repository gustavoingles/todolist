package database

import (
	"context"
	"time"
	"todo-list/user"

	"gorm.io/gorm"
)

type UserDB struct {
	ID           int64     `gorm:"primaryKey"`
	Name         string    `gorm:"not null; uniqueIndex; check: name >= 4"`
	PasswordHash string    `gorm:"column:password_hash;not null; check: passwordHash >= 6"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type UserStore struct {
	db *gorm.DB
}

func (u *UserStore) CreateUser(ctx context.Context, us user.User) error {
	err := gorm.G[UserDB](u.db).Create(ctx, &UserDB{
		Name:         us.Name,
		PasswordHash: us.PasswordHash,
		CreatedAt:    time.Now(),
	})
	return err
}

func (u *UserStore) GetUserById(ctx context.Context, uID int64) (user.User, error) {
	us, err := gorm.G[UserDB](u.db).Where("id = ?", uID).First(ctx)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:           us.ID,
		Name:         us.Name,
		PasswordHash: us.PasswordHash,
	}, nil
}

func (u *UserStore) GetUserByName(ctx context.Context, uName string) (user.User, error) {
	us, err := gorm.G[UserDB](u.db).Where("name = ?", uName).First(ctx)
	if err != nil {
		return user.User{}, err
	}

	return user.User{
		ID:           us.ID,
		Name:         us.Name,
		PasswordHash: us.PasswordHash,
	}, nil
}

func (u *UserStore) UpdateUserById(ctx context.Context, uID int64, f user.UserUpdateData) error {
	err := u.db.Transaction(func(tx *gorm.DB) error {
		user, err := gorm.G[UserDB](tx).Where("id = ?", uID).First(ctx)
		if err != nil {
			return err
		}

		updateFn := func(uDB *UserDB) UserDB {
			uDB.Name = f.NewName
			uDB.PasswordHash = f.NewPassword
			return *uDB
		}

		_, err = gorm.G[UserDB](tx).Updates(ctx, updateFn(&user))
		return err
	})
	return err
}

func (u *UserStore) DeleteUserById(ctx context.Context, uID int64) error {
	_, err := gorm.G[UserDB](u.db).Where("id = ?", uID).Delete(ctx)
	return err
}
