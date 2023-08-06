package storeservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type (
	Coordinates struct {
		lat  float64
		long float64
	}

	Store struct {
		Name        string `json:"empresa"`
		Coordinates Coordinates
		Address     string `json:"direccion"`
		Localidad   string `json:"localidad"`
		Latitud     string `json:"latitud"`
		Longitud    string `json:"longitud"`
	}

	StoreIterator interface {
		HasNext() bool
		GetNext() (*Store, error)
	}

	fetchStoreIterator struct {
		fetchUrl  string
		address   string
		localidad string
		index     int
		maxCount int
		stores    []Store
	}

	Result struct {
		Stores []Store `json:"data"`
		MaxCount int `json:"recordsFiltered"`
	}
)

func (it fetchStoreIterator) GetFetchURL() string {
	s := `{"draw":2,"columns":[{"data":"id","name":"id","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"empresa","name":"empresa","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"direccion","name":"direccion","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"localidad","name":"localidad","searchable":true,"orderable":false,"search":{"value":"%s","regex":false}},{"data":"rubro","name":"rubro","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"qr","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"id","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"latitud","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}},{"data":"longitud","name":"","searchable":true,"orderable":false,"search":{"value":"","regex":false}}],"order":[],"start":%d,"length":50,"search":{"value":"","regex":false}}`
	return fmt.Sprintf(s, it.localidad, it.index)
}

func (s *Store) UnmarshalJSON(data []byte) error {
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
	return it.index <= it.maxCount
}

func (it fetchStoreIterator) fetchNextBatch() error {
response, err := http.Get(it.GetFetchURL())

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

	return nil
}

func (it fetchStoreIterator) GetNext() (*Store, error) {

	if !it.HasNext() {
		return nil, errors.New("no next value")
	}

	if len(it.stores) > it.index {
		store := &it.stores[it.index]
		it.index++
		return store, nil
	}

	err := it.fetchNextBatch()

	if err != nil {
		return nil, err
	}

if len(it.stores) > it.index {
		store := &it.stores[it.index]
		it.index++
		return store, nil
	} else {
		return nil, errors.New("unknown error")
	}

	
}

func GetStoresIterator(fetchUrl string, address string, localidad string) StoreIterator {
	return fetchStoreIterator{fetchUrl: fetchUrl, address: address, localidad: localidad}
}
