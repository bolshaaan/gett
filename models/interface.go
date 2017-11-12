package models

/*
	Interface, that helps to make mock Driver object
	and test handlers
*/
type DriverI interface {
	BeforeSave() error
	Create() error
	GetID() uint
	Validate() error
}
