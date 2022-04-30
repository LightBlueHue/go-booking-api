package services

import (
	"go-booking-api/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IDBService interface {
	InitDB()
	GetDB() *gorm.DB
}

type DBService struct {
}

var db *gorm.DB

func GetDBService() IDBService {

	return &DBService{}
}

func (dbService *DBService) InitDB() {

	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "user=postgres password=postgres dbname=go-booking-api port=5432 sslmode=disable TimeZone=Europe/London",
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {

		panic("failed to connect database")
	}

	var dbResult *gorm.DB

	// Migrate the schema
	if tiExists := database.Migrator().HasTable(&models.TicketInventory{}); !tiExists {

		if err = database.AutoMigrate(&models.TicketInventory{}); err != nil {

			panic("ite")
		}

		if err = database.AutoMigrate(&models.Booking{}); err != nil {

			panic("bte")
		}

		if dbResult = database.Create(&models.TicketInventory{AvailableTickets: 50, TotalTickets: 50, Name: "JusticeLeagueLive", Description: "Justice League Live"}); dbResult.Error != nil {

			panic("ipe")
		}

		// if dbResult = database.Exec("DROP FUNCTION book(INT,INT)"); dbResult.Error != nil {

		// 	panic("dfb")
		// }

		if dbResult = database.Exec(`CREATE OR REPLACE FUNCTION book(ticketsToBuy INT, userId INT)
		RETURNS SETOF BIGINT AS
		$body$
			DECLARE
			remainingTickets INT;
			timeNow timestamp with time zone;
			
			BEGIN
			
			lock table ticket_inventories IN ROW EXCLUSIVE MODE;
			
			remainingTickets:= (SELECT ticket_inventories.available_tickets FROM ticket_inventories WHERE ID = 1);
			
			IF remainingTickets > ticketsToBuy THEN
				BEGIN
				
				timeNow := now();
				
				UPDATE ticket_inventories
				SET available_tickets = ( remainingTickets - ticketsToBuy)
				WHERE id = 1;
				
				RETURN QUERY INSERT INTO bookings (created_at, updated_at, user_id, tickets, ticket_inventory_id)
				VALUES (timeNow, timeNow, userId, ticketsToBuy, 1) RETURNING id;
				
				END;
			
			END IF;
			
			END 
		$body$
		LANGUAGE plpgsql`); dbResult.Error != nil {

			panic("utcbf")
		}
	}

	if err = database.AutoMigrate(&models.Credentials{}); err != nil {

		panic("cte")
	}

	if err = database.AutoMigrate(&models.User{}); err != nil {

		panic("ute")
	}

	db = database
}

func (dbService *DBService) GetDB() *gorm.DB {

	if db == nil {

		dbService.InitDB()
	}
	return db
}
