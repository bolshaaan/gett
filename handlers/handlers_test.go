package handlers

import (
	"testing"
	"github.com/valyala/fasthttp"
	"github.com/stretchr/testify/assert"
	"github.com/bolshaaan/gett/models"
	"encoding/json"
	"github.com/stretchr/testify/require"
)

func TestGetHandler(t *testing.T) {
	t.Run("No id passed", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		GetHandler( ctx )
		assert.Equal(t, ctx.Response.StatusCode(), fasthttp.StatusBadRequest)
	})

	t.Run("Invalid Driver id", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.SetUserValue("id", "-122")
		GetHandler( ctx )
		assert.Equal(t, ctx.Response.StatusCode(), fasthttp.StatusBadRequest)
	})

	t.Run("Driver exists", func(t *testing.T) {
		GetDriver = func( id uint ) (*models.Driver, error ) {
			return &models.Driver{}, nil
		}

		ctx := &fasthttp.RequestCtx{}
		ctx.SetUserValue("id", "123")
		GetHandler(ctx)

		assert.Equal(t, ctx.Response.StatusCode(), fasthttp.StatusOK )
		assert.Equal(t, ctx.Response.Header.Peek("Content-Type"), []byte("application/json") )
	})
}

func TestImportHandler(t *testing.T) {
	t.Run("Invalid json", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}
		ctx.SetBodyString( "malformed json string" )
		ImportHandler( ctx )
		assert.Equal(t, ctx.Response.StatusCode(), fasthttp.StatusBadRequest)
	})

	t.Run("Success import driver", func(t *testing.T) {

		ctx := &fasthttp.RequestCtx{}
		ctx.SetBodyString( `["license_number"]` )

		getDriversFromRequest = func(ctx *fasthttp.RequestCtx) ([]models.DriverI, error) {
			return []models.DriverI{
				models.DriverI( MockDriver{ ID: 10, Name: "Bill Gates", LicenseNumber: "1222-1333-12" } ),
			}, nil
		}

		ImportHandler( ctx )

		assert.Equal(t, ctx.Response.StatusCode(), fasthttp.StatusOK)
	})

	t.Run("Create error of some driver", func(t *testing.T) {
		ctx := &fasthttp.RequestCtx{}

		getDriversFromRequest = func(ctx *fasthttp.RequestCtx) ([]models.DriverI, error) {
			return []models.DriverI{
				models.DriverI( MockDriver{ ID: 10, Name: "Bill Gates", LicenseNumber: "1222-1333-12" } ),
				models.DriverI( MockDriver{ ID: BadDriverID, Name: "Unlucky driver", LicenseNumber: "1222-1333-77" } ),
			}, nil
		}

		ImportHandler( ctx )

		assert.Equal(t, ctx.Response.StatusCode(), fasthttp.StatusOK)
		assert.Equal(t, ctx.Response.Header.Peek("Content-Type"), []byte("application/json"))

		errReport := []ErrReport{}
		err := json.Unmarshal( ctx.Response.Body(), &errReport )
		require.NoError(t, err)

		require.True(t, len(errReport) > 0)
		assert.Equal(t, errReport[0].Id, BadDriverID)
	})
}
