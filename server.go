package gett

import (
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	log "github.com/sirupsen/logrus"
	"github.com/bolshaaan/gett/db"
	"os"
	"os/signal"
	"syscall"
	"github.com/bolshaaan/gett/models"
)

var gracefulStop = make(chan os.Signal)

func StartApp( addr, pgUrl string ) {
	db.InitDB(pgUrl)
	db.DB.AutoMigrate( &models.Driver{} ) // just creates table, if not exists
	defer db.DB.Close()

	router := fasthttprouter.New()

	router.POST("/import", ImportHandler)
	router.GET("/driver/:id", GetHandler)

	log.Infof("Starting server at %s", addr)

	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		// for gracefull stop
		<-gracefulStop

		db.DB.Close()
		os.Exit(0)
	}()

	log.Fatal(fasthttp.ListenAndServe(addr, router.Handler))
}

func StopApp() {
	gracefulStop <- os.Interrupt
}
