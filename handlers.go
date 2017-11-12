package gett

import (
	"github.com/valyala/fasthttp"
	"encoding/json"
	"github.com/bolshaaan/gett/models"
	log "github.com/sirupsen/logrus"
	"fmt"
	"strconv"
)

// able to mock method
var GetDriver = models.GetDriver

func GetHandler(ctx *fasthttp.RequestCtx) {
	strId, ok := ctx.UserValue("id").(string)
	if ! ok {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(strId, 10, 0)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	driver, err := GetDriver(uint(id))

	switch {
	case err == models.ErrNoDriver:
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	case err != nil:
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	jDriver, err := json.Marshal( driver )
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetContentType("application/json")
	fmt.Fprint(ctx, string(jDriver))

	return
}

type ErrReport struct {
	Id uint `json:"id"`
	Err string `json:"err"`
}

// parses json body from request
var getDriversFromRequest = func (ctx *fasthttp.RequestCtx) ([]models.DriverI, error) {
	drivers := []models.Driver{}

	if err := json.Unmarshal(ctx.PostBody(), &drivers); err != nil {
		return nil, err
	}

	drvs := make([]models.DriverI, len(drivers) )
	for i, v := range drivers {
		drvs[i] = models.DriverI( v )
	}

	return drvs, nil
}

// import receives body in json format
// and one by one add to database
func ImportHandler(ctx *fasthttp.RequestCtx) {
	drivers, err := getDriversFromRequest(ctx)

	if err != nil {
		log.Error("Bad json body: ", err)
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	if drivers == nil {
		return
	}
	errReport := []*ErrReport{}

	for _, dr := range drivers {
		if err := dr.Create(); err != nil {
			errReport = append(errReport, &ErrReport{ Id: dr.GetID(), Err:err.Error() } )
		}
	}

	if len(errReport) > 0 {
		ctx.SetContentType("application/json")
		text, err := json.Marshal( errReport );
		if err != nil {
			log.Error("Problems while making report: ", err)
		}
		fmt.Fprint(ctx, string(text))
	}
}
