package db

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
)

var DB *gorm.DB
//var pgScheme = `postgresql://postgres@127.0.0.1/gett?sslmode=disable`

func InitDB(pgUrl string) {

	var err error
	DB, err = gorm.Open("postgres", pgUrl)

	if err != nil {
		log.Fatal(err)
	}
}

