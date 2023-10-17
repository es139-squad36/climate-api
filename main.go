package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type solarData struct {
	Consumption float32
	Cost        float32
	Hour        string
}

type windData struct {
	Consumption float32
	Cost        float32
	Hour        string
}

type emissionsData struct {
	Emit float32
	Day  string
}

var solar = []solarData{}
var wind = []windData{}
var emissions = []emissionsData{}

func mockData() {
	currentTime := time.Now()

	// Loop through the past 7 days
	for i := 0; i < 7; i++ {
		// Calculate the date for the current day
		date := currentTime.AddDate(0, 0, -(7 - i))

		emit := 100.0 + rand.Float32()*(300.0-100.0)
		e := emissionsData{
			Emit: emit,
			Day:  date.Format(time.ANSIC),
		}
		emissions = append(emissions, e)

		// Loop through 24 hours
		for hour := 0; hour < 24; hour++ {
			// Calculate the time for the current hour
			timestamp := date.Add(time.Hour * time.Duration(hour))
			solarCons := 20.0 + rand.Float32()*(30.0-20.0)
			solarCost := 30.0 + rand.Float32()*(60.0-30.0)
			windCons := 10.0 + rand.Float32()*(40.0-10.0)
			windCost := 30.0 + rand.Float32()*(40.0-30.0)
			s := solarData{
				Consumption: solarCons,
				Cost:        solarCost,
				Hour:        timestamp.Format(time.ANSIC),
			}
			solar = append(solar, s)
			w := windData{
				Consumption: windCons,
				Cost:        windCost,
				Hour:        timestamp.Format(time.ANSIC),
			}
			wind = append(wind, w)
		}
	}
}

func getData(context *gin.Context) {
	// Convert to JSON
	data := map[string]interface{}{
		"solar":     solar,
		"wind":      wind,
		"emissions": emissions,
	}
	context.JSON(http.StatusOK, data)
}

func main() {
	// Generate mock data
	mockData()

	// Init gin
	router := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"} // Allow requests from any origin
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Setup API
	router.GET("/getData", getData)

	// Create a ticker that ticks every hour
	ticker := time.NewTicker(time.Hour)

	// Function to run every hour
	periodicFunction := func() {
		timestamp := time.Now()
		emissionsTime, err := time.Parse(time.ANSIC, emissions[len(emissions)-1].Day)
		duration := 24 * time.Hour
		if err != nil {
			duration = timestamp.Sub(emissionsTime)
		}
		if duration == 24*time.Hour {
			emit := 100.0 + rand.Float32()*(300.0-100.0)
			e := emissionsData{
				Emit: emit,
				Day:  timestamp.Format(time.ANSIC),
			}
			emissions = append(emissions[1:], e)
		}
		solarCons := 20.0 + rand.Float32()*(30.0-20.0)
		solarCost := 30.0 + rand.Float32()*(60.0-30.0)
		windCons := 10.0 + rand.Float32()*(40.0-10.0)
		windCost := 30.0 + rand.Float32()*(40.0-30.0)
		s := solarData{
			Consumption: solarCons,
			Cost:        solarCost,
			Hour:        timestamp.Format(time.ANSIC),
		}
		solar = append(solar[1:], s)
		w := windData{
			Consumption: windCons,
			Cost:        windCost,
			Hour:        timestamp.Format(time.ANSIC),
		}
		wind = append(wind[1:], w)
	}

	// Start a goroutine to execute the function at the specified interval
	go func() {
		for {
			select {
			case <-ticker.C:
				periodicFunction()
			}
		}
	}()

	router.Run()
}
