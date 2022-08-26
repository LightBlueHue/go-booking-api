package services

import (
	"fmt"
	"go-booking-api/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	SQL_STATEMENT_CREATE_DB            = `CREATE DATABASE "%s" WITH OWNER = postgres ENCODING = 'UTF8' LC_COLLATE = 'en_US.utf8' LC_CTYPE = 'en_US.utf8' TABLESPACE = pg_default CONNECTION LIMIT = -1;`
	SQL_STATEMENT_CREATE_BOOK_FUNCTION = `CREATE OR REPLACE FUNCTION book(ticketsToBuy INT, userId INT)
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
	LANGUAGE plpgsql`
)

type DBService struct {
	db *gorm.DB
}

type DbInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	SslMode  string
	TimeZone string
}

func NewDBService(db *gorm.DB) IDBService {

	return &DBService{db}
}

// InitDB initialises database and creates tables. It will by create a new database if one does not exist.
func (s *DBService) InitDB(dbInfo DbInfo, open DbInitialiser, createDbStatement string) *gorm.DB {

	var dbResult *gorm.DB
	var err error
	connStrWithDbName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", dbInfo.Host, dbInfo.User, dbInfo.Password, dbInfo.DbName, dbInfo.Port, dbInfo.SslMode, dbInfo.TimeZone)
	connStrWithoutDbName := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=%s TimeZone=%s", dbInfo.Host, dbInfo.User, dbInfo.Password, dbInfo.Port, dbInfo.SslMode, dbInfo.TimeZone)

	s.db, err = open(postgres.New(postgres.Config{
		DSN:                  connStrWithDbName,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {

		s.db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  connStrWithoutDbName,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})

		if dbResult = s.db.Exec(createDbStatement); dbResult.Error != nil {

			panic("failed to connect database")
		}

		s.db, err = gorm.Open(postgres.New(postgres.Config{
			DSN:                  connStrWithDbName,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
	}

	// Migrate the schema
	if tiExists := s.db.Migrator().HasTable(&models.TicketInventory{}); !tiExists {

		if err = s.db.AutoMigrate(&models.TicketInventory{}); err != nil {

			panic("ite")
		}

		if err = s.db.AutoMigrate(&models.Booking{}); err != nil {

			panic("bte")
		}

		if dbResult = s.db.Create(&models.TicketInventory{AvailableTickets: 50, TotalTickets: 50, Name: "JusticeLeagueLive", Description: "Justice League Live"}); dbResult.Error != nil {

			panic("ipe")
		}

		if dbResult = s.db.Exec(SQL_STATEMENT_CREATE_BOOK_FUNCTION); dbResult.Error != nil {

			panic("utcbf")
		}
	}

	if err = s.db.AutoMigrate(&models.Credentials{}); err != nil {

		panic("cte")
	}

	if err = s.db.AutoMigrate(&models.User{}); err != nil {

		panic("ute")
	}

	return s.db
}
