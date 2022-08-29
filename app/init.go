package app

import (
	"fmt"
	"go-booking-api/app/controllers"
	"go-booking-api/app/services"

	_ "github.com/revel/modules"
	"github.com/revel/revel"
	"gorm.io/gorm"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	db *gorm.DB
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		ServicesFilter,
		//revel.SessionFilter,     // Restore and write the session cookie.
		//revel.FlashFilter,       // Restore and write the flash cookie.
		revel.ValidationFilter,  // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,        // Resolve the requested language
		HeaderFilter,            // Add some security based headers
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.CompressFilter,    // Compress the result.
		revel.BeforeAfterFilter, // Call the before and after filter functions
		revel.ActionInvoker,     // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

var ServicesFilter = func(c *revel.Controller, fc []revel.Filter) {

	if db == nil {

		initDB()
	}

	if ac, ok := c.AppController.(*controllers.AccountController); ok && !ac.Service.IsServiceInitialized() {

		ac.Service = services.NewService(db)
	}

	if bc, ok := c.AppController.(*controllers.BookingController); ok && !bc.Service.IsServiceInitialized() {

		bc.Service = services.NewService(db)
	}

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}

func initDB() {

	dbInfo := newDBInfo()
	db = services.NewDBService(db).InitDB(dbInfo, gorm.Open, fmt.Sprintf(services.SQL_STATEMENT_CREATE_DB, dbInfo.DbName))

	if db == nil {
		panic("INIT DB")
	}
}

func newDBInfo() services.DbInfo {

	dbInfo := services.DbInfo{}

	dbInfo.Host = revel.Config.StringDefault("db.host", "")
	dbInfo.DbName = revel.Config.StringDefault("db.dbname", "")
	dbInfo.Password = revel.Config.StringDefault("db.password", "")
	dbInfo.Port = revel.Config.IntDefault("db.port", 0)
	dbInfo.SslMode = revel.Config.StringDefault("db.sslmode", "")
	dbInfo.TimeZone = revel.Config.StringDefault("db.tz", "")
	dbInfo.User = revel.Config.StringDefault("db.user", "")
	return dbInfo
}
