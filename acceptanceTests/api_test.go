package acceptanceTests

import (
	"testing"
	"github.com/bolshaaan/gett"
	"time"
	"net/http"
	"strings"
	"github.com/magiconair/properties/assert"
	"fmt"
	dbPq "github.com/bolshaaan/gett/db"
	"encoding/json"
	"github.com/bolshaaan/gett/models"
)

var testAddr = "127.0.0.1:8081"

// !!! MUST BE TEST DATABASE, all data will be removed !!!
var testPgScheme = `postgresql://postgres@127.0.0.1/gett?sslmode=disable`

func prepareTestServer() {
	dbPq.InitDB(testPgScheme)
	dbPq.DB.DropTableIfExists("drivers")
	dbPq.DB.AutoMigrate()

	go func() {
		gett.StartApp( testAddr,testPgScheme)
	}()

	time.Sleep(time.Second)
}

func gracefulStop() {
	 dbPq.DB.Exec("truncate table drivers")
	 gett.StopApp()
}

func TestAPI(t *testing.T) {
	prepareTestServer()
	defer gracefulStop()

	importUrl := fmt.Sprintf("http://%s/import", testAddr)
	getUrl := fmt.Sprintf("http://%s/driver/", testAddr)

	t.Run("Create 3 drivers + 1 repeat id", func(t *testing.T) {
		drivers := []models.Driver{
			{ ID: 1, Name: "Johny b. Goode", LicenseNumber: "12-234-45" },
			{ ID: 2, Name: "Taylor Swift", LicenseNumber: "12-234-10" },
			{ ID: 333, Name: "Eyal Golan", LicenseNumber: "12-288-10" },

			// driver with repeat id
			{ ID: 1, Name: "Johnatan", LicenseNumber: "12-11-222" },
		}

		input, _ := json.Marshal(drivers)
		resp, err := http.Post(importUrl, "application/json",  strings.NewReader(string(input)) )
		if err != nil {
			t.Fatal("Error: ", err)
		}

		dec := json.NewDecoder(resp.Body)

		report := []gett.ErrReport{}
		if err := dec.Decode(&report); err != nil {
			t.Fatal("error decoding ", err)
		}

		assert.Equal(t, resp.StatusCode, http.StatusOK)
		assert.Equal(t, len(report), 1, "1 bad element in bulk operation")

		t.Run("Get 1 exist driver", func(t *testing.T) {
			resp, err := http.Get(getUrl + "1")
			if err != nil {
				t.Fatal("Error: ", err)
			}

			assert.Equal(t, resp.StatusCode, http.StatusOK)

			dec := json.NewDecoder(resp.Body)

			dr := &models.Driver{}
			if err := dec.Decode( dr ); err != nil {
				t.Fatal("error of decoding: ", err)
			}

			assert.Equal(t, dr.Name, drivers[0].Name)
			assert.Equal(t, dr.LicenseNumber, drivers[0].LicenseNumber)
		})
	})
}
