package services

import (
	"errors"
	"go-booking-api/app/models"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"syreclabs.com/go/faker"
)

type mockJWTService struct {
	mock.Mock
}

func Test_EmailExists_Returns_True(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	email := faker.Internet().Email()

	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM users WHERE email = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(strconv.Itoa(1)))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actual := target.EmailExists(email)

	assert.True(t, actual)
}

func Test_EmailExists_Returns_False(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	email := faker.Internet().Email()

	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(id) FROM users WHERE email = $1")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).
			AddRow(strconv.Itoa(0)))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actual := target.EmailExists(email)

	assert.False(t, actual)
}

func Test_GetPassword_Returns_Data(t *testing.T) {

	var db *gorm.DB
	var actualError error
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	pwd := faker.Internet().Password(9, 20)

	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ")).
		WithArgs(pwd).
		WillReturnRows(sqlmock.NewRows([]string{"password"}).
			AddRow(pwd))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actual, actualError := target.GetPassword(pwd)

	assert.Nil(t, actualError)
	assert.Equal(t, pwd, actual)
}

func Test_GetPassword_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var actualError error
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	expectedError := errors.New("my error")
	pwd := faker.Internet().Password(9, 20)

	sqlMock.ExpectQuery(regexp.QuoteMeta("SELECT password FROM users INNER JOIN credentials ON users.credential_id = credentials.id WHERE email = ")).
		WithArgs(pwd).
		WillReturnError(expectedError)

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actualPwd, actualError := target.GetPassword(pwd)

	assert.Empty(t, actualPwd)
	assert.Error(t, actualError)
	assert.Equal(t, expectedError, actualError)
}

func Test_GetByEmail_Returns_User(t *testing.T) {

	var db *gorm.DB
	var setupError error
	var defaultTime time.Time

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	utcNow := time.Now().UTC()
	email := faker.Internet().Email()

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, "yac", "yaccadamia", "yac@yac.co", "1"))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actualUser, actualError := target.GetByEmail(email)

	assert.Nil(t, actualError)
	assert.NotNil(t, actualUser)
}

func Test_GetByEmail_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var setupError error
	expectedError := errors.New("my error")

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
	target := GetUserService(db)

	actualUser, actualError := target.GetByEmail(email)

	assert.Equal(t, expectedError, actualError)
	assert.Error(t, actualError)
	assert.Nil(t, actualUser)
}

func Test_GetByToken_Returns_User(t *testing.T) {

	var db *gorm.DB
	var setupError error
	var defaultTime time.Time

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	var jwts = newMockJWTService(t)
	utcNow := time.Now().UTC()
	token := faker.RandomString(20)
	email := faker.Internet().Email()
	jwts.On("GetClaim", token, EMAIL_CLAIM).Return(email, nil)

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, "yac", "yaccadamia", "yac@yac.co", "1"))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actualUser, actualError := target.GetByToken(token, jwts)

	assert.Nil(t, actualError)
	assert.NotNil(t, actualUser)
}

func Test_GetByToken_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var setupError error
	var defaultTime time.Time

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	var jwts = newMockJWTService(t)
	utcNow := time.Now().UTC()
	token := faker.RandomString(20)
	expectedError := errors.New("my error")
	email := ""
	jwts.On("GetClaim", token, EMAIL_CLAIM).Return(email, expectedError)

	sqlMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, "yac", "yaccadamia", "yac@yac.co", "1"))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actualUser, actualError := target.GetByToken(token, jwts)

	assert.Equal(t, expectedError, actualError)
	assert.Error(t, actualError)
	assert.Nil(t, actualUser)
}

func Test_Save_Returns_NoError(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	user := &models.User{FirstName: "yac", LastName: "yaccadamia", CredentialID: 1}
	itIsAny := sqlmock.AnyArg()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","first_name","last_name","email","credential_id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow("1"))
	sqlMock.ExpectCommit()

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actualError := target.Save(user)

	assert.Nil(t, actualError)
}

func Test_Save_Returns_Error(t *testing.T) {

	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	user := &models.User{FirstName: "yac", LastName: "yaccadamia", CredentialID: 1}
	itIsAny := sqlmock.AnyArg()
	expectedError := errors.New("my error")

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users" ("created_at","updated_at","deleted_at","first_name","last_name","email","credential_id") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "id"`)).
		WithArgs(itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny).
		WillReturnError(expectedError)
	sqlMock.ExpectRollback()

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := GetUserService(db)

	actualError := target.Save(user)

	assert.Error(t, actualError)
	assert.Equal(t, expectedError, actualError)
}

func (m *mockJWTService) GetClaim(token string, claimType JwtClaimType) (string, error) {
	ret := m.Called(token, claimType)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, JwtClaimType) string); ok {
		r0 = rf(token, claimType)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, JwtClaimType) error); ok {
		r1 = rf(token, claimType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m *mockJWTService) GenerateToken(email string, isUser bool) string {

	return ""
}

func (m *mockJWTService) ValidateToken(token string) (*jwt.Token, error) {

	return nil, nil
}

func newMockJWTService(t testing.TB) *mockJWTService {

	mock := &mockJWTService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
