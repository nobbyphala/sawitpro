package repository

import (
	"context"
	"database/sql"
	"errors"
	"sawitpro/entity"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_userProfileRepository_UpdateProfileById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx        context.Context
		tx         *sqlx.Tx
		id         string
		updateData entity.UserProfile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		mock    func()
	}{
		{
			name: "successfully update profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				id:  "id1",
				updateData: entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62867",
				},
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectExec("UPDATE user_profile").WithArgs("jonathan", "+62867", "id1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "successfully update profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				id: "id1",
				updateData: entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62867",
				},
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectExec("UPDATE user_profile").WithArgs("jonathan", "+62867", "id1").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "error update profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				id:  "id1",
				updateData: entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62867",
				},
			},
			wantErr: errors.New("error update"),
			mock: func() {
				mock.ExpectExec("UPDATE user_profile").WithArgs("jonathan", "+62867", "id1").WillReturnError(errors.New("error update"))
			},
		},
		{
			name: "error update profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				id: "id1",
				updateData: entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62867",
				},
			},
			wantErr: errors.New("error update"),
			mock: func() {
				mock.ExpectExec("UPDATE user_profile").WithArgs("jonathan", "+62867", "id1").WillReturnError(errors.New("error update"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			repo := userProfileRepository{
				db: tt.fields.db,
			}
			err := repo.UpdateProfileById(tt.args.ctx, tt.args.tx, tt.args.id, tt.args.updateData)
			assert.Equal(t, tt.wantErr, err)
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestNewUserProfileRepository(t *testing.T) {
	db, _, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type args struct {
		db *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want userProfileRepository
	}{
		{
			name: "return instance of user profile repository",
			args: args{
				db: dbx,
			},
			want: userProfileRepository{
				db: dbx,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserProfileRepository(tt.args.db)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_userProfileRepository_InsertProfile(t *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx  context.Context
		tx   *sqlx.Tx
		user entity.UserProfile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr error
		mock    func()
	}{
		{
			name: "success insert profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				user: entity.UserProfile{
					FullName:    "phala",
					PhoneNumber: "+627876234",
					Password:    "l;asuidfjuahsaso;idjf",
				},
			},
			want:    "profile_id_1",
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("INSERT INTO user_profile").
					WithArgs("phala", "+627876234", "l;asuidfjuahsaso;idjf").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("profile_id_1"))
			},
		},
		{
			name: "success insert profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				user: entity.UserProfile{
					FullName:    "phala",
					PhoneNumber: "+627876234",
					Password:    "l;asuidfjuahsaso;idjf",
				},
			},
			want:    "profile_id_1",
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("INSERT INTO user_profile").
					WithArgs("phala", "+627876234", "l;asuidfjuahsaso;idjf").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("profile_id_1"))
			},
		},
		{
			name: "error insert profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				user: entity.UserProfile{
					FullName:    "phala",
					PhoneNumber: "+627876234",
					Password:    "l;asuidfjuahsaso;idjf",
				},
			},
			want:    "",
			wantErr: errors.New("error insert"),
			mock: func() {
				mock.ExpectQuery("INSERT INTO user_profile").
					WithArgs("phala", "+627876234", "l;asuidfjuahsaso;idjf").
					WillReturnError(errors.New("error insert"))
			},
		},
		{
			name: "success insert profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				user: entity.UserProfile{
					FullName:    "phala",
					PhoneNumber: "+627876234",
					Password:    "l;asuidfjuahsaso;idjf",
				},
			},
			want:    "",
			wantErr: errors.New("error insert"),
			mock: func() {
				mock.ExpectQuery("INSERT INTO user_profile").
					WithArgs("phala", "+627876234", "l;asuidfjuahsaso;idjf").
					WillReturnError(errors.New("error insert"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := userProfileRepository{
				db: tt.fields.db,
			}
			got, err := repo.InsertProfile(tt.args.ctx, tt.args.tx, tt.args.user)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_userProfileRepository_GetProfileById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx context.Context
		tx  *sqlx.Tx
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.UserProfile
		wantErr error
		mock    func()
	}{
		{
			name: "success get profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				id:  "profile-id-1",
			},
			want: entity.UserProfile{
				Id:          "profile-id-1",
				FullName:    "phala",
				PhoneNumber: "+621234",
				Password:    "qwerty",
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"profile-id-1",
				).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"full_name",
						"phone_number",
						"password",
					}).AddRow(
						"profile-id-1",
						"phala",
						"+621234",
						"qwerty",
					),
				)
			},
		},
		{
			name: "success get profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				id: "profile-id-1",
			},
			want: entity.UserProfile{
				Id:          "profile-id-1",
				FullName:    "phala",
				PhoneNumber: "+621234",
				Password:    "qwerty",
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"profile-id-1",
				).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"full_name",
						"phone_number",
						"password",
					}).AddRow(
						"profile-id-1",
						"phala",
						"+621234",
						"qwerty",
					),
				)
			},
		},
		{
			name: "profile not found",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				id:  "profile-id-1",
			},
			want:    entity.UserProfile{},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"profile-id-1",
				).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "profile not found with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				id: "profile-id-1",
			},
			want:    entity.UserProfile{},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"profile-id-1",
				).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "error get profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx:  nil,
				id:  "profile-id-1",
			},
			want:    entity.UserProfile{},
			wantErr: errors.New("error select"),
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"profile-id-1",
				).WillReturnError(errors.New("error select"))
			},
		},
		{
			name: "error get profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				id: "profile-id-1",
			},
			want:    entity.UserProfile{},
			wantErr: errors.New("error select"),
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"profile-id-1",
				).WillReturnError(errors.New("error select"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := userProfileRepository{
				db: tt.fields.db,
			}
			got, err := repo.GetProfileById(tt.args.ctx, tt.args.tx, tt.args.id)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_userProfileRepository_GetProfileByPhoneNumber(t *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx         context.Context
		tx          *sqlx.Tx
		phoneNumber string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.UserProfile
		wantErr error
		mock    func()
	}{
		{
			name: "success get profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx:         context.TODO(),
				tx:          nil,
				phoneNumber: "+621234",
			},
			want: entity.UserProfile{
				Id:          "profile-id-1",
				FullName:    "phala",
				PhoneNumber: "+621234",
				Password:    "qwerty",
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"+621234",
				).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"full_name",
						"phone_number",
						"password",
					}).AddRow(
						"profile-id-1",
						"phala",
						"+621234",
						"qwerty",
					),
				)
			},
		},
		{
			name: "success get profile with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				phoneNumber: "+621234",
			},
			want: entity.UserProfile{
				Id:          "profile-id-1",
				FullName:    "phala",
				PhoneNumber: "+621234",
				Password:    "qwerty",
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"+621234",
				).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"full_name",
						"phone_number",
						"password",
					}).AddRow(
						"profile-id-1",
						"phala",
						"+621234",
						"qwerty",
					),
				)
			},
		},
		{
			name: "profile not found",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx:         context.TODO(),
				tx:          nil,
				phoneNumber: "+621234",
			},
			want:    entity.UserProfile{},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"+621234",
				).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "profile not found with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				phoneNumber: "+621234",
			},
			want:    entity.UserProfile{},
			wantErr: nil,
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"+621234",
				).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "error get profile",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx:         context.TODO(),
				tx:          nil,
				phoneNumber: "+621234",
			},
			want:    entity.UserProfile{},
			wantErr: errors.New("error select"),
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"+621234",
				).WillReturnError(errors.New("error select"))
			},
		},
		{
			name: "profile not found with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				phoneNumber: "+621234",
			},
			want:    entity.UserProfile{},
			wantErr: errors.New("error select"),
			mock: func() {
				mock.ExpectQuery("SELECT").WithArgs(
					"+621234",
				).WillReturnError(errors.New("error select"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := userProfileRepository{
				db: tt.fields.db,
			}
			got, err := repo.GetProfileByPhoneNumber(tt.args.ctx, tt.args.tx, tt.args.phoneNumber)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_userProfileRepository_IncreaseSuccessLoginCount(t *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx       context.Context
		tx        *sqlx.Tx
		profileId string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		mock    func()
	}{
		{
			name: "success update",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx:       context.TODO(),
				tx:        nil,
				profileId: "profile-id-1",
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectExec("UPDATE").WithArgs(
					"profile-id-1",
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "success update with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				profileId: "profile-id-1",
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectExec("UPDATE").WithArgs(
					"profile-id-1",
				).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "got error when update",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx:       context.TODO(),
				tx:        nil,
				profileId: "profile-id-1",
			},
			wantErr: errors.New("error update"),
			mock: func() {
				mock.ExpectExec("UPDATE").WithArgs(
					"profile-id-1",
				).WillReturnError(errors.New("error update"))
			},
		},
		{
			name: "got update error with transaction",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				tx: func() *sqlx.Tx {
					mock.ExpectBegin()
					tx, _ := dbx.Beginx()
					return tx
				}(),
				profileId: "profile-id-1",
			},
			wantErr: errors.New("error update"),
			mock: func() {
				mock.ExpectExec("UPDATE").WithArgs(
					"profile-id-1",
				).WillReturnError(errors.New("error update"))
			},
		},
	}
	for _, tt := range tests {
		tt.mock()

		t.Run(tt.name, func(t *testing.T) {
			repo := userProfileRepository{
				db: tt.fields.db,
			}
			err := repo.IncreaseSuccessLoginCount(tt.args.ctx, tt.args.tx, tt.args.profileId)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_userProfileRepository_RunWithTransaction(t *testing.T) {
	db, mock, _ := sqlmock.New()
	dbx := sqlx.NewDb(db, "pgx")

	type fields struct {
		db *sqlx.DB
	}
	type args struct {
		ctx        context.Context
		handleFunc func(tx *sqlx.Tx) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		mock    func()
	}{
		{
			name: "handler function success",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				handleFunc: func(tx *sqlx.Tx) error {
					return nil
				},
			},
			wantErr: nil,
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
		},
		{
			name: "handler function error",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				handleFunc: func(tx *sqlx.Tx) error {
					return errors.New("error handler")
				},
			},
			wantErr: errors.New("error handler"),
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
		},
		{
			name: "commit error",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				handleFunc: func(tx *sqlx.Tx) error {
					return nil
				},
			},
			wantErr: errors.New("error commit"),
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectCommit().WillReturnError(errors.New("error commit"))
			},
		},
		{
			name: "transaction begin error",
			fields: fields{
				db: dbx,
			},
			args: args{
				ctx: context.TODO(),
				handleFunc: func(tx *sqlx.Tx) error {
					return nil
				},
			},
			wantErr: errors.New("error begin"),
			mock: func() {
				mock.ExpectBegin().WillReturnError(errors.New("error begin"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			repo := userProfileRepository{
				db: tt.fields.db,
			}
			err := repo.RunWithTransaction(tt.args.ctx, tt.args.handleFunc)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
