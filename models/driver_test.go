package models

import (
	"testing"
	dbPq "github.com/bolshaaan/gett/db"
	"github.com/jinzhu/gorm"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/pkg/errors"
	"github.com/lib/pq"
)

func Setup() sqlmock.Sqlmock{
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("Can't create sqlmock %s", err)
	}

	dbPq.DB, err = gorm.Open("postgres", db)
	if err != nil {
		log.Fatalf("Cant init gorm connection %s", err)
	}

	return mock
}

func TestDriver_Create(t *testing.T) {
	mock := Setup()

	dr := &Driver{ID: 123, Name: "Kolya", LicenseNumber: "license-number"}

	t.Run("Success create", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id"}).
			AddRow("123")

		mock.ExpectQuery(`INSERT INTO "drivers"`).WillReturnRows(rows)

		err := dr.Create()
		require.NoError(t, err)

		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

	t.Run("Unique primary key", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO "drivers"`).
			WillReturnError( &pq.Error{ Constraint: "drivers_pkey" } )

		err := dr.Create()
		require.Error(t, err)
		assert.Equal(t, ErrDriverIDExists, err)

		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

	t.Run("Unique license", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO "drivers"`).
			WillReturnError( &pq.Error{ Constraint: "drivers_license_number_key" } )

		err := dr.Create()
		require.Error(t, err)
		assert.Equal(t, ErrUniqLicense, err)

		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

	t.Run("Some db error", func(t *testing.T) {
		mock.ExpectQuery(`INSERT INTO "drivers"`).
			WillReturnError( &pq.Error{ Message: "Something awful"  } )

		err := dr.Create()
		require.Error(t, err)
		assert.Equal(t, ErrDBError, err)

		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

	t.Run("No name", func(t *testing.T) {
		dr := &Driver{ Name: "", LicenseNumber: "1231" }

		err := dr.Create()
		require.Error(t, err)
	})

	t.Run("No license number", func(t *testing.T) {
		dr := &Driver{ Name: "", LicenseNumber: "1231" }

		err := dr.Create()
		require.Error(t, err)
	})
}

func TestGetDriver(t *testing.T) {
	mock := Setup()

	t.Run("Driver exists", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "name", "license_number"}).
			AddRow("123", "Pizhon", "123-33-11")

		mock.ExpectQuery(`SELECT \* FROM "drivers"`).
			WithArgs(123).
			WillReturnRows(rows)

		driver, err := GetDriver( 123 )
		require.NoError(t, err)

		assert.Equal(t, uint(123), driver.ID )
		assert.Equal(t, "Pizhon", driver.Name )

		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

	t.Run("Driver NOT exists", func(t *testing.T) {
		rows := sqlmock.
			NewRows([]string{"id", "name", "license_number"})

		mock.ExpectQuery(`SELECT \* FROM "drivers"`).
			WithArgs(123).
			WillReturnRows(rows)

		_, err := GetDriver( 123 )
		require.Error(t, err)

		assert.Equal(t, ErrNoDriver, err)
		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

	t.Run("DB error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM "drivers"`).
			WithArgs(123).WillReturnError( errors.New("some db error") )

		_, err := GetDriver( 123 )
		require.Error(t, err)

		assert.Equal(t, ErrDBError, err)
		assert.NoError(t, mock.ExpectationsWereMet(), "Check expectations")
	})

}

