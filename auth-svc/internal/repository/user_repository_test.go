package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cloud9cloud9/go-grpc-todo/auth-svc/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

const (
	testEmail         = "test@example.com"
	hashedPassword    = "hashed_password"
	testErrorEmail    = "test_error@example.com"
	testNotFoundEmail = "test_not_found@example.com"
)

var (
	args = []string{"id", "email", "password"}
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestAuthPostgres_FindByEmail(t *testing.T) {
	gormDB, mock, err := setupMockDB()
	assert.NoError(t, err)

	defer func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}()

	userRepo := NewRepository(gormDB)

	tests := []struct {
		name          string
		email         string
		mockSetup     func()
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name:  "User Found",
			email: testEmail,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT (1|\$2)`).
					WithArgs(testEmail, 1).
					WillReturnRows(sqlmock.NewRows(args).
						AddRow(0, testEmail, hashedPassword))
			},
			expectedUser:  &domain.User{Email: testEmail, Password: hashedPassword},
			expectedError: nil,
		},
		{
			name:  "User Not Found",
			email: testNotFoundEmail,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT (1|\$2)`).
					WithArgs(testNotFoundEmail, 1).
					WillReturnRows(sqlmock.NewRows(args))
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name:  "Database Error",
			email: testErrorEmail,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"."id" LIMIT (1|\$2)`).
					WithArgs(testErrorEmail, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := userRepo.FindByEmail(tt.email)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedUser, user)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthPostgres_FindByID(t *testing.T) {
	gormDB, mock, err := setupMockDB()
	assert.NoError(t, err)

	defer func() {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
	}()

	userRepo := NewRepository(gormDB)

	tests := []struct {
		name          string
		id            int64
		mockSetup     func()
		expectedUser  *domain.User
		expectedError error
	}{
		{
			name: "User Found",
			id:   1,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"."id" LIMIT (1|\$2)`).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows(args).
						AddRow(1, testEmail, hashedPassword))
			},
			expectedUser:  &domain.User{Id: 1, Email: testEmail, Password: hashedPassword},
			expectedError: nil,
		},
		{
			name: "User Not Found",
			id:   2,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"."id" LIMIT (1|\$2)`).
					WithArgs(2, 1).
					WillReturnRows(sqlmock.NewRows(args))
			},
			expectedUser:  nil,
			expectedError: gorm.ErrRecordNotFound,
		},
		{
			name: "Database Error",
			id:   3,
			mockSetup: func() {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE id = \$1 ORDER BY "users"."id" LIMIT (1|\$2)`).
					WithArgs(3, 1).
					WillReturnError(sql.ErrConnDone)
			},
			expectedUser:  nil,
			expectedError: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			user, err := userRepo.FindByID(tt.id)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedUser, user)
		})
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
