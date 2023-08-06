package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type (
	GeoCodeStats struct {
		Latitude  string `json:"lat"`
		Longitude string `json:"lon"`
	}
)

const apiUrl = "https://geocode.maps.co/search"

func (s *Place) SetCoordinates() {
	c := s.Coordinates
	if c.lat == 0 || c.long == 0 {
		time.Sleep(1 * time.Second)
		params := url.Values{}
		params.Add("q", s.Address+" "+s.Localidad)
		response, err := http.Get(apiUrl + "?" + params.Encode())

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		defer response.Body.Close()

		// Read the response body
		data, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		var result []GeoCodeStats
		if err := json.Unmarshal(data, &result); err != nil {
			fmt.Println("Error:", err)
			return
		}

		if len(result) > 0 {
			s.Coordinates.lat, _ = strconv.ParseFloat(result[0].Latitude, 64)
			s.Coordinates.long, _ = strconv.ParseFloat(result[0].Longitude, 64)
		}

	}
}
