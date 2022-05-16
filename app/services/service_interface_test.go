package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_SetServices_Updates_Correctly(t *testing.T) {

	os.Setenv(GO_BOOKING_API_SECRET, "E59DD115760893782F7FB8CC6C387DE86FFEC3C186A8EFE24184E9CABDB2EFC3")
	var db *gorm.DB
	dbs := GetDBService(db)
	hs := GetHashService()
	jwts := GetJWTService()
	rs := GetResponseService()
	us := GetUserService(db)
	vs := GetValidationService()
	bs := GetBookingService(db)
	target := Service{}

	target.SetServices(dbs, hs, jwts, rs, us, vs, bs)

	assert.Equal(t, dbs, target.DBService)
	assert.Equal(t, hs, target.HashService)
	assert.Equal(t, jwts, target.JWTService)
	assert.Equal(t, rs, target.ResponseService)
	assert.Equal(t, us, target.UserService)
	assert.Equal(t, vs, target.ValidationService)
	assert.Equal(t, bs, target.BookingService)
}
