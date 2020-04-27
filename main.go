package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/BenJuan26/elite"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type event struct {
	Mod *string `json:"mod,omitempty"`
	Key string  `json:"key"`
}

var events []event

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	var ok bool
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		err = fmt.Errorf("createEvent: Empty body, %s", err.Error())
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	if err = json.Unmarshal(reqBody, &newEvent); err != nil {
		err = fmt.Errorf("Unmarshal error: %s\n%s", err.Error(), string(reqBody))
		log.Errorf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	log.Debugf("Event received: %+v", newEvent)

	if newEvent.Mod == nil { //single keypress
		if ok, err = SendKeyPress(newEvent.Key); !ok {
			err = fmt.Errorf("createEvent: SendKeyPress error %s", err.Error())
			log.Errorf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
	} else { // keypress with modified
		if ok, err = SendInput(false, *newEvent.Mod); !ok {
			err = fmt.Errorf("createEvent: SendInput error %s", err.Error())
			log.Errorf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		time.Sleep(10 * time.Millisecond)

		if ok, err = SendKeyPress(newEvent.Key); !ok {
			err = fmt.Errorf("createEvent: SendKeyPress error %s", err.Error())
			log.Errorf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}

		time.Sleep(10 * time.Millisecond)

		if ok, err = SendInput(true, *newEvent.Mod); !ok {
			err = fmt.Errorf("createEvent: SendInput error %s", err.Error())
			log.Errorf(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEvent)
}

func edSystem(w http.ResponseWriter, r *http.Request) {

	var err error
	system, _ := elite.GetStarSystem()
	jsonEnc := json.NewEncoder(w)
	err = jsonEnc.Encode(system)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func edStatus(w http.ResponseWriter, r *http.Request) {

	var err error
	var status *elite.Status

	if status, err = elite.GetStatus(); err != nil {
		log.Error(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	/*status.Flags.LandingGearDown = true
	status.Flags.FlightAssistOff = false
	status.Flags.LightsOn = true
	status.Flags.CargoScoopDeployed = true*/

	log.Debugf("Status: %+v", status)
	if err = json.NewEncoder(w).Encode(status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type StatusRespWr struct {
	http.ResponseWriter // We embed http.ResponseWriter
	status              int
}

func (w *StatusRespWr) WriteHeader(status int) {
	w.status = status // Store the status for our own use
	w.ResponseWriter.WriteHeader(status)
}
func wrapHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		srw := &StatusRespWr{ResponseWriter: w}
		h.ServeHTTP(srw, r)
		if srw.status >= 400 { // 400+ codes are the error codes
			log.Errorf("Error status code: %d when serving path: %s",
				srw.status, r.RequestURI)
		}
	}
}

const (
	STATIC_DIR = "/html/"
	PORT       = "8080"
)

func main() {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Info("Starting EDVBB v0.2")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/system", edSystem).Methods("GET")
	router.HandleFunc("/status", edStatus).Methods("GET")
	router.PathPrefix(STATIC_DIR).Handler(
		http.StripPrefix(
			STATIC_DIR,
			handlers.
				LoggingHandler(os.Stdout, http.
					FileServer(http.Dir("."+STATIC_DIR)))))
	log.Fatal(http.ListenAndServe(":"+PORT, router))
}
