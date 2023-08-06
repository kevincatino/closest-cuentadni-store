package main

import (
	"closest-cuentadni-store/services"
	"fmt"
	"os"
	"sort"

	"github.com/joho/godotenv"
)

type Places struct {
	stores    []*services.Place
	reference services.Coordinates
}

// Implement the Sort interface for ByAge
func (a Places) Len() int { return len(a.stores) }
func (a Places) Less(i, j int) bool {
	return a.stores[i].Coordinates.GetDistance(a.reference) < a.stores[j].Coordinates.GetDistance(a.reference)
}
func (a Places) Swap(i, j int) {
	a.stores[i], a.stores[j] = a.stores[j], a.stores[i]
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	localidad := os.Getenv("LOCALIDAD")
	address := os.Getenv("ADDRESS")
	var place = services.Place{
		Localidad: localidad,
		Address:   address,
	}
	place.SetCoordinates()
	url := os.Getenv("URL")
	apiUrl, err := services.GetAPIUrl(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	iterator := services.GetStoresIterator(apiUrl, localidad)
	var stores []services.Place = iterator.ToArray()
	var storesMap map[string]*services.Place = make(map[string]*services.Place)
	var uniqueStores []*services.Place
	for idx, store := range stores {
		key := store.Address + " " + store.Localidad
		if _, ok := storesMap[key]; ok {
			fmt.Println("store repeated")
			continue
		}
		store.SetCoordinates()
		uniqueStores = append(uniqueStores, &stores[idx]) // cannot use &store because the same variable is being reused in the loop so the pointer value would be the same for all values of the array
		storesMap[key] = &stores[idx]
		fmt.Printf("%+v\n", store)
	}

	sort.Sort(Places(Places{stores: uniqueStores, reference: place.Coordinates}))

	fmt.Println("\n\n\n\nClosest stores:")
	for i := 0; i < len(uniqueStores) && i < 20; i++ {
		fmt.Printf("%s, %s, %s\n", uniqueStores[i].Name, uniqueStores[i].Address, uniqueStores[i].Localidad)
	}

}
