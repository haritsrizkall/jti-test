package mysql

import (
	"context"
	"database/sql"

	"github.com/haritsrizkall/jti-test/domain"
)

type mysqlRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &mysqlRepository{
		db: db,
	}
}

func (r *mysqlRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	row := r.db.QueryRowContext(ctx, GET_BY_EMAIL, email)

	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *mysqlRepository) Store(ctx context.Context, user domain.User) (int, error) {
	result, err := r.db.ExecContext(ctx, STORE, user.Name, user.Email, user.Password)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
