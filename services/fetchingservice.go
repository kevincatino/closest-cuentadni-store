package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

const baseURL = "https://www.bancoprovincia.com.ar/cuentadni/Home/GetLocalesListadoByIdBuscador?idBuscador="

func GetAPIUrl(url string) (string, error) {
	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	defer response.Body.Close()

	// Read the response body
	data, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	input := string(data)
	pattern := `idBuscador=[\+\s']+(\d+),`

	// Compile the regex pattern
	re := regexp.MustCompile(pattern)

	// Find all matches of the pattern in the input string
	match := re.FindStringSubmatch(input)

	// Extract the number from the matches
	if len(match) >= 2 {
		number, err := strconv.Atoi(match[1])
		if err != nil {
			fmt.Println("Error converting to integer:", err)
			return "", err
		}
		return fmt.Sprintf("%s%d",baseURL,number), nil
	} else {
		fmt.Println("Number not found.")
		return "", errors.New("number not found")
	}
}
