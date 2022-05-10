package tests

import (
	"errors"
	"go-booking-api/app/services"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"syreclabs.com/go/faker"
)

func Test_EmailExists_ReturnsTrue(t *testing.T) {

	var db *gorm.DB
	var err error
	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	email := faker.Internet().Email()
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM users WHERE email = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(strconv.Itoa(1)))

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, err)
	target := services.GetUserService(db)

	actual := target.EmailExists(email)

	assert.True(t, actual)
}

func Test_EmailExists_ReturnsFalse(t *testing.T) {

	var db *gorm.DB
	var err error
	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	email := faker.Internet().Email()
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM users WHERE email = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(strconv.Itoa(0)))

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, err)
	target := services.GetUserService(db)

	actual := target.EmailExists(email)

	assert.False(t, actual)
}

func Test_GetPassword_ReturnsData(t *testing.T) {

	var db *gorm.DB
	var actualError error
	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	pwd := faker.Internet().Password(9, 20)
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ")).
		WithArgs(pwd).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).
			AddRow(pwd))

	db, actualError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, actualError)
	target := services.GetUserService(db)

	actual, actualError := target.GetPassword(pwd)

	assert.Nil(t, actualError)
	assert.Equal(t, pwd, actual)
}

func Test_GetPassword_ReturnsError(t *testing.T) {

	var db *gorm.DB
	var actualError error
	var setupError error
	expectedError := errors.New("mock error")
	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	pwd := faker.Internet().Password(9, 20)
	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ")).
		WithArgs(pwd).
		WillReturnError(expectedError)

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.GetUserService(db)

	actualPwd, actualError := target.GetPassword(pwd)

	assert.Empty(t, actualPwd)
	assert.Equal(t, expectedError, actualError)
}

func Test_GetByEmail_ReturnsUser(t *testing.T) {

	var db *gorm.DB
	var err error
	var setupError error
	var defaultTime time.Time
	utcNow := time.Now().UTC()
	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	email := faker.Internet().Email()
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, "yac", "yaccadamia", "yac@yac.co", "1"))

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.GetUserService(db)

	actualUser, err := target.GetByEmail(email)

	assert.Nil(t, err)
	assert.NotEmpty(t, actualUser)
}

func Test_GetByEmail_ReturnsError(t *testing.T) {

	var db *gorm.DB
	var setupError error
	expectedError := errors.New("Mock error")

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	email := faker.Internet().Email()
	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(expectedError)
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.GetUserService(db)

	actualUser, actualError := target.GetByEmail(email)

	assert.Equal(t, expectedError, actualError)
	assert.Empty(t, actualUser)
}
