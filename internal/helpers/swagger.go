package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

var swaggerMap = sync.Map{}

const (
	swaggerKey         = "swagger"
	defaultSwaggerPath = "api/generated/app.swagger.json"
)

// ServeSwagger - serves swagger specification.
func ServeSwagger(w http.ResponseWriter, r *http.Request) {
	_, ok := swaggerMap.Load(swaggerKey)
	if !ok {
		swaggerPath := defaultSwaggerPath
		if overridedPath, ok := os.LookupEnv("SWAGGER_PATH"); ok {
			swaggerPath = overridedPath
		}
		m, err := getSwaggerMap(swaggerPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		swaggerMap.Store(swaggerKey, m)
	}
	// serving a swagger file
	m, _ := swaggerMap.Load("swagger")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(m)
	if err != nil {
		fmt.Println(err)
	}
}

func getSwaggerMap(swaggerFileName string) (map[string]interface{}, error) {
	// read a swagger json file
	jsonFile, err := os.Open(swaggerFileName)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// map a json file to the struct
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(byteValue, &m)
	if err != nil {
		return nil, err
	}
	// return a result
	return m, nil
}
