package main

import (
	"math"
	"os"
)

type (
	Coordinates struct {
		lat  float64
		long float64
	}

	Store struct {
		name        string
		coordinates Coordinates
		address     string
		localidad   string
	}
)

func (c Coordinates) Lat() float64 {
	return c.lat
}

func (c Coordinates) Long() float64 {
	return c.long
}

func (c Coordinates) GetDistance(other Coordinates) float64 {
	return math.Sqrt(math.Pow(c.Lat()-other.Lat(), 2) + math.Pow(c.Long()-other.Long(), 2))
}

func main() {
	localidad := os.Getenv("LOCALIDAD")
	address := os.Getenv("DIRECCION")
}
