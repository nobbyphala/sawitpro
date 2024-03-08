package repository

import (
	"context"
	"sawitpro/entity"

	"github.com/jmoiron/sqlx"
)

type TransactionHandleFunc = func(tx *sqlx.Tx) error

type UserProfileRepositoryInterface interface {
	RunWithTransaction(ctx context.Context, handleFunc TransactionHandleFunc) error
	InsertProfile(ctx context.Context, tx *sqlx.Tx, user entity.UserProfile) (string, error)
	GetProfileById(ctx context.Context, tx *sqlx.Tx, id string) (entity.UserProfile, error)
	UpdateProfileById(ctx context.Context, tx *sqlx.Tx, id string, updateData entity.UserProfile) error
	GetProfileByPhoneNumber(ctx context.Context, tx *sqlx.Tx, phoneNumber string) (entity.UserProfile, error)
	IncreaseSuccessLoginCount(ctx context.Context, tx *sqlx.Tx, profileId string) error
}
