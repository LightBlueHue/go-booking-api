package services_test

import (
	"errors"
	"go-booking-api/app/models"
	"go-booking-api/app/services"
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

func Test_EmailExists_WhenEmailExistsInDb_Returns_True(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	email := faker.Internet().Email()

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(strconv.Itoa(1)))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actual := target.EmailExists(email)

	// Assert
	assert.True(t, actual)
}

func Test_EmailExists_WhenEmailDoesNotExistInDb_Returns_False(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	email := faker.Internet().Email()

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(strconv.Itoa(0)))

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actual := target.EmailExists(email)

	// Assert
	assert.False(t, actual)
}

func Test_GetPassword_WhenPasswordExistsInDb_Returns_Password(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var actualError error
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	pwd := faker.Internet().Password(9, 20)

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(pwd).
		WillReturnRows(sqlmock.NewRows([]string{""}).
			AddRow(pwd))

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actual, actualError := target.GetPassword(pwd)

	// Assert
	assert.Nil(t, actualError)
	assert.Equal(t, pwd, actual)
}

func Test_GetPassword_WhenDbReturnsError_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var actualError error
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	expectedError := errors.New("my error")
	pwd := faker.Internet().Password(9, 20)

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(pwd).
		WillReturnError(expectedError)

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualPwd, actualError := target.GetPassword(pwd)

	// Assert
	assert.Empty(t, actualPwd)
	assert.Error(t, actualError)
	assert.Equal(t, expectedError, actualError)
}

func Test_GetByEmail_WhenEmailExistsInDb_Returns_User(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error
	var defaultTime time.Time

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()

	utcNow := time.Now().UTC()
	email := faker.Internet().Email()

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, "yac", "yaccadamia", "yac@yac.co", "1"))

	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualUser, actualError := target.GetByEmail(email)

	// Assert
	assert.Nil(t, actualError)
	assert.NotNil(t, actualUser)
}

func Test_GetByEmail_WhenEmailDoesNotExistInDb_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error
	expectedError := errors.New("my error")

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	email := faker.Internet().Email()
	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(email).
		WillReturnError(expectedError)

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualUser, actualError := target.GetByEmail(email)

	// Assert
	assert.Equal(t, expectedError, actualError)
	assert.Error(t, actualError)
	assert.Nil(t, actualUser)
}

func Test_GetByToken_WhenUserExistsInDb_Returns_User(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error
	var defaultTime time.Time

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	var jwts = newMockJWTService(t)
	utcNow := time.Now().UTC()
	token := faker.RandomString(20)
	expectedUser := &models.User{FirstName: "yac", LastName: "yaccadamia", Email: "yac@yac.co", CredentialID: 1}
	jwts.On("GetClaim", token, services.EMAIL_CLAIM).Return(expectedUser.Email, nil)

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(expectedUser.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, expectedUser.FirstName, expectedUser.LastName, expectedUser.Email, expectedUser.CredentialID))

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualUser, actualError := target.GetByToken(token, jwts)

	// Assert
	assert.Nil(t, actualError)
	assert.NotNil(t, actualUser)
	assert.Equal(t, expectedUser.FirstName, actualUser.FirstName)
	assert.Equal(t, expectedUser.LastName, actualUser.LastName)
	assert.Equal(t, expectedUser.Email, actualUser.Email)
	assert.Equal(t, expectedUser.CredentialID, actualUser.CredentialID)
}

func Test_GetByToken_WhenClaimInvalid_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error
	var defaultTime time.Time

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	var jwts = newMockJWTService(t)
	utcNow := time.Now().UTC()
	token := faker.RandomString(20)
	expectedError := errors.New("my error")
	email := ""
	jwts.On("GetClaim", token, services.EMAIL_CLAIM).Return(email, expectedError)

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email", "credential_id"}).
			AddRow("1", utcNow, utcNow, defaultTime, "yac", "yaccadamia", "yac@yac.co", "1"))

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualUser, actualError := target.GetByToken(token, jwts)

	// Assert
	assert.Equal(t, expectedError, actualError)
	assert.Error(t, actualError)
	assert.Nil(t, actualUser)
}

func Test_GetByToken_WhenUserDoesNotExistInDb_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	var jwts = newMockJWTService(t)
	token := faker.RandomString(20)
	email := faker.Internet().Email()
	jwts.On("GetClaim", token, services.EMAIL_CLAIM).Return(email, nil)
	expectedError := errors.New("my error")

	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(email).
		WillReturnError(expectedError)

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualUser, actualError := target.GetByToken(token, jwts)

	// Assert
	assert.Equal(t, expectedError, actualError)
	assert.Error(t, actualError)
	assert.Nil(t, actualUser)
}

func Test_Save_WhenNoError_Returns_NoError(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	user := &models.User{FirstName: "yac", LastName: "yaccadamia", CredentialID: 1}
	itIsAny := sqlmock.AnyArg()

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow("1"))
	sqlMock.ExpectCommit()

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualError := target.Save(user)

	// Assert
	assert.Nil(t, actualError)
}

func Test_Save_WhenDbError_Returns_Error(t *testing.T) {

	// Arrange
	var db *gorm.DB
	var setupError error

	sqlDb, sqlMock, _ := sqlmock.New()
	defer sqlDb.Close()
	db, setupError = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	user := &models.User{FirstName: "yac", LastName: "yaccadamia", CredentialID: 1}
	itIsAny := sqlmock.AnyArg()
	expectedError := errors.New("my error")

	sqlMock.ExpectBegin()
	sqlMock.ExpectQuery(regexp.QuoteMeta("")).
		WithArgs(itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny, itIsAny).
		WillReturnError(expectedError)
	sqlMock.ExpectRollback()

	assert.Nil(t, setupError)
	target := services.NewUserService(db)

	// Act
	actualError := target.Save(user)

	// Assert
	assert.Error(t, actualError)
	assert.Equal(t, expectedError, actualError)
}

func (m *mockJWTService) GetClaim(token string, claimType services.JwtClaimType) (string, error) {
	ret := m.Called(token, claimType)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, services.JwtClaimType) string); ok {
		r0 = rf(token, claimType)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, services.JwtClaimType) error); ok {
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
