package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
)

type (
	Coordinates struct {
		lat  float64
		long float64
	}

	Place struct {
		Name        string `json:"empresa"`
		Coordinates Coordinates
		Address     string `json:"direccion"`
		Localidad   string `json:"localidad"`
	}

	StoreIterator interface {
		HasNext() bool
		GetNext() (*Place, error)
		ToArray() []Place
	}

	fetchStoreIterator struct {
		fetchUrl  string
		localidad string
		index     int
		maxCount  int
		stores    []Place
	}

	Result struct {
		Stores   []Place `json:"data"`
		MaxCount int     `json:"recordsFiltered"`
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

func (it fetchStoreIterator) GetFetchBody() string {
	body := `{"draw":2,"columns":[{"data":"id","name":"id","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"empresa","name":"empresa","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"direccion","name":"direccion","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"localidad","name":"localidad","searchable":true,"orderable":false,"search":{"value":"%s","regex":false}},{"data":"rubro","name":"rubro","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"qr","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"id","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"latitud","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"longitud","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[],"start":%d,"length":50,"search":{"value":"","regex":false}}`
	return fmt.Sprintf(body, it.localidad, it.index)
}

func (it fetchStoreIterator) ToArray() []Place {
	for it.maxCount > len(it.stores) {
		err := it.fetchNextBatch()
		if err != nil {
			fmt.Println(err.Error())
			return it.stores
		}
	}
	return it.stores
}

func (s *Place) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to hold the flat JSON data
	type FlatStore struct {
		Name      string  `json:"empresa"`
		Address   string  `json:"direccion"`
		Localidad string  `json:"localidad"`
		Latitud   float64 `json:"latitud"`
		Longitud  float64 `json:"longitud"`
	}

	// Unmarshal the flat JSON data into the temporary struct
	var flatStore FlatStore
	if err := json.Unmarshal(data, &flatStore); err != nil {
		return err
	}

	// Assign the values to the nested struct fields
	s.Address = flatStore.Address
	s.Coordinates = Coordinates{
		lat:  flatStore.Latitud,
		long: flatStore.Longitud,
	}
	s.Localidad = flatStore.Localidad
	s.Name = flatStore.Name

	return nil
}

func (it fetchStoreIterator) HasNext() bool {
	return it.index < it.maxCount
}

func (it *fetchStoreIterator) fetchNextBatch() error {
	response, err := http.Post(it.fetchUrl, "application/json", strings.NewReader(it.GetFetchBody()))

	if err != nil {
		return err
	}

	defer response.Body.Close()

	// Read the response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	var result Result

	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	it.stores = append(it.stores, result.Stores...)
	it.maxCount = result.MaxCount

	return nil
}

func (it *fetchStoreIterator) GetNext() (*Place, error) {

	if !it.HasNext() {
		return nil, errors.New("no next value")
	}

	if it.index < len(it.stores) {
		store := &it.stores[it.index]
		it.index++
		return store, nil
	}

	err := it.fetchNextBatch()

	if err != nil {
		return nil, err
	}

	if it.index < len(it.stores) {
		store := &it.stores[it.index]
		it.index++
		return store, nil
	} else {
		return nil, errors.New("unknown error")
	}

}

func GetStoresIterator(fetchUrl string, localidad string) StoreIterator {
	it := fetchStoreIterator{fetchUrl: fetchUrl, localidad: localidad}
	it.fetchNextBatch()
	return &it
}
