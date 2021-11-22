package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	_ "asengoz/yemekSepetiRestApi/swagger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//KeyValue struct of value godoc
type KeyValue struct {
	Value string `json:"value"`
}

//FlushData struct of all values godoc
type FlushData struct {
	Values []KeyValue `json:"values"`
}

//key-value variable godoc
type singleton map[string]string

var (
	instance singleton
	filename = "keyValuesJSON.json"
)

//getInstance return instance godoc
func getInstance() singleton {

	if instance == nil {

		instance = make(singleton)
	}

	return instance
}

//checkFile checks for the existence of file. If it doesn't exist, it creates a new file godoc
func checkFile(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

//writeJSON write data into filename json godoc
func writeJSON(filename string) error{
	flushResult := readDataFromJSON()

	newValue := &KeyValue{
		Value: instance["key"],
	}

	flushResult.Values = append(flushResult.Values, *newValue)

	// Preparing the data to be marshalled and written.
	dataBytes, writeErr := json.Marshal(flushResult)
	if writeErr != nil {
		fmt.Println(writeErr)
		return writeErr
	}

	writeErr = ioutil.WriteFile(filename, dataBytes, 0644)
	if writeErr != nil {
		fmt.Println(writeErr)
		return writeErr
	}
	return nil
}

//writeJSONBackgroundTask runs at a certain time interval and writes the value to the file godoc
func writeJSONBackgroundTask() {
	ticker := time.NewTicker(1 * time.Minute)
	for _ = range ticker.C {
		writeJSON(filename)
	}
}

//readDataFromJSON read data from JSON godoc
func readDataFromJSON() FlushData {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var flushResult FlushData
	json.Unmarshal(byteValue, &flushResult)
	return flushResult
}

//Get get the value from key godoc
// @Summary gets the value from key
// @Description gets the value from key
// @Tags methods
// @Accept */*
// @Produce json
// @Success 200 {string} string "success"
// @Router /get [get]
func Get(c *gin.Context) {
	if instance["key"] == "" {
		instance["key"] = "DefaultValue"
	}

	value := instance["key"]

	c.String(http.StatusOK, value)
}

//Post set the value into key godoc
// @Summary set the value into key
// @Description set the value into key
// @Tags methods
// @Accept json
// @Param value body KeyValue true "value"
// @Produce json
// @Success 200 {string} string "success"
// @Router /post [post]
func Post(c *gin.Context) {
	/*
		Example of Post JSON Body
		{
		    "value": "test"
		}
	*/
	var keyValue KeyValue

	err := c.BindJSON(&keyValue)
	if err != nil {
		log.Fatal(err)
	}
	instance["key"] = keyValue.Value

	c.String(http.StatusOK, "New value of key is %s", instance["key"])
}

//Flush get all data from JSON godoc
// @Summary gets all data
// @Description gets all data from JSON
// @Tags methods
// @Accept */*
// @Produce json
// @Success 200 {string} string "success"
// @Router /flush [get]
func Flush(c *gin.Context) {
	flushResult := readDataFromJSON()
	for i := 0; i < len(flushResult.Values); i++ {
		c.String(http.StatusOK, "%d. value is %s \n", i+1, flushResult.Values[i])
	}
}

// @title Yemeksepeti Golang
// @version 1.0
// @description Atakan Şengöz - Yemeksepeti Api Swagger
// @termsOfService http://swagger.io/terms/
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:9090
// @BasePath /
// @schemes http
func main() {

	//check file
	err := checkFile(filename)
	if err != nil {
		fmt.Println(err)
	}

	//get value
	getInstance()

	//value is saved in a specified interval into file
	go writeJSONBackgroundTask()

	//when the program is opened, it reads if there is data in the file and calls its last element.
	temp := readDataFromJSON()
	if  temp.Values != nil && len(temp.Values)>0 {
		instance["key"] = temp.Values[len(temp.Values)-1].Value
	}

	//create an instance of Gin Framework
	r := gin.Default()

	//A new route for SwaggerUI
	url := ginSwagger.URL("http://localhost:9090/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//an endpoint to get the key
	r.GET("/get", Get)

	//an endpoint to set the key. it needs a json to set key
	r.POST("/post", Post)

	//an endpoint to flush the all data
	r.GET("/flush", Flush)

	//listen and serve on localhost:9090
	r.Run(":9090")

}
