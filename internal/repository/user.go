package repository

import (
	"context"
	"database/sql"
	"golang-yt/domain"

	"github.com/doug-martin/goqu/v9"
)

type userRepository struct {
	db *goqu.Database
}

func NewUser(con *sql.DB) domain.UserRepository {
	return &userRepository{
		db: goqu.New("postgres", con),
	}
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (user domain.User, err error) {
	dataset := u.db.From("users").Where(goqu.C("email").Eq(email))
	_, err = dataset.ScanStructContext(ctx, &user)
	return
}
