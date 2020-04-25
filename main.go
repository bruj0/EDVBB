package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/BenJuan26/elite"
	"github.com/gorilla/mux"
	"github.com/micmonay/keybd_event"
)

type event struct {
	Key int `json:"key"`
}

var events []event
var kb keybd_event.KeyBonding

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Error(w, "Empty body")
	}

	if err := json.Unmarshal(reqBody, &newEvent); err != nil {
		log.Errorf("Unmarshal error: %s", string(reqBody))
		return
	}

	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	log.Debugf("Event received: %+v", newEvent)

	kb.SetKeys(newEvent.Key)

	err = kb.Launching()
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(newEvent)
}

func edSystem(w http.ResponseWriter, r *http.Request) {

	var err error
	system, _ := elite.GetStarSystem()

	if err = json.NewEncoder(w).Encode(system); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func edStatus(w http.ResponseWriter, r *http.Request) {

	var err error
	status, _ := elite.GetStatus()

	if err = json.NewEncoder(w).Encode(status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func main() {

	var err error

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Info("Starting remoto v0.1")

	kb, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/system", edSystem).Methods("GET")
	router.HandleFunc("/status", edStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
