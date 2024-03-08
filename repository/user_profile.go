package repository

import (
	"context"
	"database/sql"
	"sawitpro/entity"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type userProfileRepository struct {
	db *sqlx.DB
}

func NewUserProfileRepository(db *sqlx.DB) userProfileRepository {
	return userProfileRepository{
		db: db,
	}
}

func (repo userProfileRepository) InsertProfile(ctx context.Context, tx *sqlx.Tx, user entity.UserProfile) (string, error) {
	var id string
	var err error

	if tx != nil {
		err = tx.QueryRowContext(
			ctx,
			queryInserProfile,
			user.FullName,
			user.PhoneNumber,
			user.Password,
		).Scan(&id)
	} else {
		err = repo.db.QueryRowContext(
			ctx,
			queryInserProfile,
			user.FullName,
			user.PhoneNumber,
			user.Password,
		).Scan(&id)
	}

	return id, err
}

func (repo userProfileRepository) GetProfileById(ctx context.Context, tx *sqlx.Tx, id string) (entity.UserProfile, error) {
	var res entity.UserProfile
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &res, queryGetProfileById, id)
	} else {
		err = repo.db.GetContext(ctx, &res, queryGetProfileById, id)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}

		return res, err
	}

	return res, nil
}

func (repo userProfileRepository) GetProfileByPhoneNumber(ctx context.Context, tx *sqlx.Tx, phoneNumber string) (entity.UserProfile, error) {
	var res entity.UserProfile
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &res, queryGetProfileByPhoneNumber, phoneNumber)
	} else {
		err = repo.db.GetContext(ctx, &res, queryGetProfileByPhoneNumber, phoneNumber)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return res, nil
		}

		return res, err
	}

	return res, nil
}

func (repo userProfileRepository) UpdateProfileById(ctx context.Context, tx *sqlx.Tx, id string, updateData entity.UserProfile) error {
	var err error

	if tx != nil {
		_, err = tx.ExecContext(
			ctx,
			queryUpdateProfileById,
			updateData.FullName,
			updateData.PhoneNumber,
			id,
		)
	} else {
		_, err = repo.db.ExecContext(
			ctx,
			queryUpdateProfileById,
			updateData.FullName,
			updateData.PhoneNumber,
			id,
		)
	}

	return err
}

func (repo userProfileRepository) IncreaseSuccessLoginCount(ctx context.Context, tx *sqlx.Tx, profileId string) error {
	var err error

	if tx != nil {
		_, err = tx.ExecContext(
			ctx,
			queryIncreaseSuccessLoginCount,
			profileId,
		)
	} else {
		_, err = repo.db.ExecContext(
			ctx,
			queryIncreaseSuccessLoginCount,
			profileId,
		)
	}

	return err
}

func (repo userProfileRepository) RunWithTransaction(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) error {
	tx, err := repo.db.Beginx()
	if err != nil {
		return err
	}

	err = handleFunc(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
