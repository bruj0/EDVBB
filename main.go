package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"syscall"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"

	"github.com/BenJuan26/elite"
	"github.com/gorilla/mux"
)

type event struct {
	Key int `json:"key"`
}

var events []event

var dll = syscall.NewLazyDLL("user32.dll")
var sendInputProc = dll.NewProc("SendInput")

// static void dummy(void) { }
type keyboardInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uint64
}

type input struct {
	inputType uint32
	ki        keyboardInput
	padding   uint64
}

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

	SendInput(false)
	time.Sleep(10 * time.Millisecond)
	SendInput(true)
	//kb.SetKeys(newEvent.Key)

	// err = kb.Launching()
	// if err != nil {
	// 	panic(err)
	// }

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
	var status *elite.Status

	if status, err = elite.GetStatus(); err != nil {
		log.Error(w, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debugf("Status: %+v", status)
	if err = json.NewEncoder(w).Encode(status); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func SendInput(up bool) {
	var i input
	i.inputType = 1     //INPUT_KEYBOARD
	i.ki.wScan = 0x041E // virtual key code for a
	i.ki.wVk = 0x41
	if up == true {
		i.ki.dwFlags = 0x0002
	} else {
		i.ki.dwFlags = 0
	}
	ret, _, err := sendInputProc.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&i)),
		uintptr(unsafe.Sizeof(i)),
	)
	log.Printf("ret: %v error: %v", ret, err)

}
func main() {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
	log.Info("Starting remoto v0.2")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/system", edSystem).Methods("GET")
	router.HandleFunc("/status", edStatus).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
