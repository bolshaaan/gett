package handlers

import (
	"time"
	"github.com/pkg/errors"
)

var BadDriverID uint = 100
var mockDB = make(map[uint]*MockDriver)

type MockDriver struct {
	ID            uint   `json:"id" gorm:"primary_key"`
	Name          string `json:"name" gorm:"type:varchar(255);not null"`
	LicenseNumber string `json:"license_number" gorm:"type:varchar(255);not null;unique"`

	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (d MockDriver) Create() error {
	/* nothing do -- no database */

	if d.ID == BadDriverID {
		return errors.New("some unlucky driver")
	}

	mockDB[ d.ID ] = &d
	return nil
}
func (d MockDriver) GetID() uint {
	return d.ID
}
func (d MockDriver) BeforeSave() error {
	return nil
}
func (d MockDriver) Validate() error {
	return nil
}
