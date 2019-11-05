package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Event struct {
	gorm.Model
	Description string `json:"description"`
	Priority    uint   `json:"priority"`
}

type Result struct {
	Result      string `json:"result"`
	Description string `json:"description"`
}

func getEventByPriority(w http.ResponseWriter, r *http.Request) {
	DB, err := gorm.Open("postgres", "host=localhost port=5432 user=witcher dbname=witcher")
	defer DB.Close()

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	pr, err := strconv.Atoi(params["priority"])
	if err != nil {
		var result Result
		result.Result = "not ok"
		result.Description = "nil parameters"
		json.NewEncoder(w).Encode(result)
		return
	}
	//var event Event
	//DB.First(&event, "priority = ?", pr)
	var events []Event
	DB.Find(&events, "priority = ?", pr)
	/*if event.ID == 0 {
		var result Result
		result.Result = "not ok"
		result.Description = "nil parameters"
		json.NewEncoder(w).Encode(result)
		return
	}*/
	json.NewEncoder(w).Encode(&events)
}

func getEventById(w http.ResponseWriter, r *http.Request) {
	DB, err := gorm.Open("postgres", "host=localhost port=5432 user=witcher dbname=witcher")
	defer DB.Close()

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		var result Result
		result.Result = "not ok!"
		result.Description = "nil parameters"
		json.NewEncoder(w).Encode(result)
		return
	}
	var event Event
	DB.First(&event, "id = ?", id)

	if event.ID == 0 {
		var result Result
		result.Result = "not ok!!"
		result.Description = "nil parameters"
		json.NewEncoder(w).Encode(result)
		return
	}
	json.NewEncoder(w).Encode(&event)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	DB, _ := gorm.Open("postgres", "host=localhost port=5432 user=witcher dbname=witcher")
	defer DB.Close()

	w.Header().Set("Content-Type", "application/json")
	var event Event
	err_ := json.NewDecoder(r.Body).Decode(&event)
	if err_ != nil {
		var result Result
		result.Result = "not ok"
		result.Description = "nil parameters"
		json.NewEncoder(w).Encode(result)
		return
	}
	DB.Create(&event)
	json.NewEncoder(w).Encode(event)
}

func main() {
	DB, err := gorm.Open("postgres", "host=localhost port=5432 user=witcher dbname=witcher")
	defer DB.Close()

	if err != nil {
		fmt.Println(err)
		return
	}
	// Migrate the schema
	DB.AutoMigrate(&Event{})

	r := mux.NewRouter()
	r.HandleFunc("/event/priority/{priority}", getEventByPriority).Methods("GET")
	r.HandleFunc("/event/{id}", getEventById).Methods("GET")
	r.HandleFunc("/event", createEvent).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", r))
	// Create
	//db.Create(&Event{Description: "Broken lavochka", Priority: 1000})

	// Read
	//var product Event
	//db.First(&product, 1)                   // find product with id 1
	//db.First(&product, "description = ?", "Broken lavochka") // find product with code l1212

	//fmt.Println(product)
	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	//db.Delete(&product)
}
