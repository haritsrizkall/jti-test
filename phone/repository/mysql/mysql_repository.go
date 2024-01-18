package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/haritsrizkall/jti-test/constant"
	"github.com/haritsrizkall/jti-test/domain"
)

type phoneRepository struct {
	db *sql.DB
}

func NewPhoneRepository(db *sql.DB) domain.PhoneRepository {
	return &phoneRepository{
		db: db,
	}
}

func (r *phoneRepository) GetAll(ctx context.Context) ([]domain.Phone, error) {
	rows, err := r.db.QueryContext(ctx, GET_ALL)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var phones []domain.Phone
	for rows.Next() {
		var phone domain.Phone
		err := rows.Scan(&phone.ID, &phone.Number, &phone.Provider)
		if err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}

	return phones, nil
}

func (r *phoneRepository) GetByID(ctx context.Context, id int) (domain.Phone, error) {
	row := r.db.QueryRowContext(ctx, GET_BY_ID, id)

	var phone domain.Phone
	err := row.Scan(&phone.ID, &phone.Number, &phone.Provider)
	if err != nil {
		return domain.Phone{}, err
	}

	return phone, nil
}

func (r *phoneRepository) Update(ctx context.Context, phone domain.Phone) error {
	result, err := r.db.ExecContext(ctx, UPDATE, phone.Number, phone.Provider, phone.ID)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		return errors.New(constant.ErrNotFound)
	}

	return nil
}

func (r *phoneRepository) Store(ctx context.Context, phone domain.Phone) (int, error) {
	result, err := r.db.ExecContext(ctx, STORE, phone.Number, phone.Provider)
	if err != nil {
		return -1, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(lastInsertID), nil
}

func (r *phoneRepository) Delete(ctx context.Context, id int) error {
	result, err := r.db.ExecContext(ctx, DELETE, id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		return errors.New(constant.ErrNotFound)
	}

	return nil
}
