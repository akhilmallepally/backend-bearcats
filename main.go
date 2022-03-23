package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type EventSummary struct {
	ID        int
	Name      string
	Location  string
	Organizer string
	Date      string
}

type JsonResponse struct {
	Type    string         `json:"type"`
	Data    []EventSummary `json:"data"`
	Message string         `json:"message"`
}

func setupDB() *sql.DB {
	envMap, mapErr := godotenv.Read(".env")
	if mapErr != nil {
		fmt.Printf("Error loading .env into map[string]string\n")
		os.Exit(1)
	}

	fmt.Printf("DB USER %s\n", envMap["DB_USER"])
	fmt.Printf("DB PASS %s\n", envMap["DB_PASSWORD"])

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", envMap["DB_USER"], envMap["DB_PASSWORD"], envMap["DB_NAME"])
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

// Function for handling messages
func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

// Function for handling errors
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func getEvents(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	rows, err := db.Query(`select * from events`)
	printMessage("getting events...")
	checkErr(err)
	var eves []EventSummary

	for rows.Next() {
		var id int
		var name string
		var location string
		var organizer string
		var date string

		err = rows.Scan(&id, &name, &location, &organizer, &date)

		checkErr(err)

		eves = append(eves, EventSummary{Name: name, Location: location, Organizer: organizer, Date: date})
	}
	var response = JsonResponse{Type: "success", Data: eves}

	json.NewEncoder(w).Encode(response)
}

func DeleteOneEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	eventId := params["id"]

	var response = JsonResponse{}

	if eventId == "" {
		response = JsonResponse{Type: "error", Message: "You are missing eventId parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting event from DB")

		_, err := db.Exec("DELETE FROM events where id = $1", eventId)

		// check errors
		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The event has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func DeleteEvents(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all events...")

	_, err := db.Exec("DELETE FROM events")

	// check errors
	checkErr(err)

	printMessage("All events have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All events have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Home page")
}
func main() {
	// Init the mux router

	router := mux.NewRouter()

	// Route handles & endpoints
	router.HandleFunc("/", home)
	router.HandleFunc("/events", getEvents).Methods("GET")
	router.HandleFunc("/events/{id}", DeleteOneEvent).Methods("DELETE")
	router.HandleFunc("/events", DeleteEvents).Methods("DELETE")
	// serve the app
	fmt.Println("Server at 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
