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
	EventId   string //`json:"eventid"`
	Name      string //`json:"name"`
	Location  string //`json:"location"`
	Organizer string //`json:"organizer"`
	Date      string //`json:"date"`
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
		var eventid string
		var name string
		var location string
		var organizer string
		var date string

		err = rows.Scan(&eventid, &name, &location, &organizer, &date)

		checkErr(err)

		eves = append(eves, EventSummary{EventId: eventid, Name: name, Location: location, Organizer: organizer, Date: date})
	}
	var response = JsonResponse{Type: "success", Data: eves}

	json.NewEncoder(w).Encode(response)
}

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	eventId := r.FormValue("eventid")
	eventName := r.FormValue("name")
	locationName := r.FormValue("location")
	organizerName := r.FormValue("organizer")
	dateVal := r.FormValue("date")

	var response = JsonResponse{}

	db := setupDB()

	printMessage("Inserting movie into DB")

	fmt.Println("Inserting new movie with ID: " + eventId + " event name" + eventName + " location name " + locationName + "organizer name" + organizerName + " with date " + dateVal)

	var lastInsertID int
	err := db.QueryRow(`INSERT INTO events(eventid, name, location, organizer, date) VALUES($1, $2, $3, $4, $5);`, eventId, eventName, locationName, organizerName, dateVal).Scan(&lastInsertID)

	// check errors
	checkErr(err)

	response = JsonResponse{Type: "success", Message: "The movie has been inserted successfully!"}

	json.NewEncoder(w).Encode(response)
}

func DeleteOneEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	eventId := params["eventid"]

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

func Edit(w http.ResponseWriter, r *http.Request) {
	db := setupDB()
	params := mux.Vars(r)

	eventId := params["eventid"]
	rows, err := db.Query("SELECT * FROM Employee WHERE id=?", eventId)
	if err != nil {
		panic(err.Error())
	}
	eves := EventSummary{}
	for rows.Next() {
		var eventid string
		var name string
		var location string
		var organizer string
		var date string

		err = rows.Scan(&eventid, &name, &location, &organizer, &date)

		checkErr(err)
		eves.EventId = eventid
		eves.Name = name
		eves.Location = location
		eves.Organizer = organizer
		eves.Date = date
		//eves = append(EventSummary{EventId: eventid, Name: name, Location: location, Organizer: organizer, Date: date})
	}

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
	router.HandleFunc("/events", CreateEvent).Methods("POST")
	router.HandleFunc("/events/{id}", DeleteOneEvent).Methods("DELETE")
	router.HandleFunc("/events", DeleteEvents).Methods("DELETE")
	// serve the app
	fmt.Println("Server at 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
