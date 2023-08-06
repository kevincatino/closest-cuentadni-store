package main

import (
	"closest-cuentadni-store/services"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)



func main() {
err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file:", err)
        return
    }

	localidad := os.Getenv("LOCALIDAD")
	// address := os.Getenv("DIRECCION")
	url := os.Getenv("URL")
	iterator := services.GetStoresIterator(url, localidad)
	for iterator.HasNext() {
		next, _ := iterator.GetNext()
		fmt.Printf("%+v\n", next)
	}
	
}
