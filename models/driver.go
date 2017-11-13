package models

/*
postgres-based data module for
 */

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/lib/pq"
	dbPq  "github.com/bolshaaan/gett/db"
	"time"
)

var ErrNoDriver = errors.New("dr: driver not exists")
var ErrDBError = errors.New("dr: db error")
var ErrUniqLicense = errors.New("dr: driver with same license already exists")
var ErrDriverIDExists = errors.New("dr: driver with same id already exists")

type Driver struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Name          string `json:"name" gorm:"type:varchar(255);not null"`
	LicenseNumber string `json:"license_number" gorm:"type:varchar(255);not null;unique"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// save to databse, check validness
func (d Driver) Create() error {
	err := dbPq.DB.Create(&d).Error

	if err != nil {
		pqErr, ok := err.(*pq.Error);

		switch {
			case ok && pqErr.Constraint == "drivers_license_number_key":
				return ErrUniqLicense
			case ok && pqErr.Constraint == "drivers_pkey":
				return ErrDriverIDExists
		}

		if ok {
			log.Error("pq error: ", pqErr.Message)
		} else {
			log.Error("error: ", err)
		}

		return ErrDBError
	}

	return nil
}

func (d Driver) GetID() uint {
	return d.ID
}

// will be called from gorm
func (d Driver) BeforeSave() error {
	return d.Validate()
}

func GetDriver(id uint)(*Driver, error) {
	driver := &Driver{}
	res := dbPq.DB.First(driver, id)

	switch {
	case res.RecordNotFound():
		return nil, ErrNoDriver
	case res.Error != nil:
		log.Error("DB error: ", res.Error)
		return nil, ErrDBError
	}

	return driver,nil
}

// primitive check (all unique checks are made by database itself)
func (d Driver) Validate() error {

	if d.Name == "" {
		return errors.New("empty name")
	}

	if d.LicenseNumber == "" {
		return errors.New("LicenseNumber")
	}

	return nil
}

