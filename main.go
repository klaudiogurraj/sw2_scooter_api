package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"geo/geo"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

//go:embed escooters
var JSONFILE []byte

type EScooter struct {
	ID    string     `json:"id"`
	GPS   *geo.Point `json:"gps"`
	KM    int64      `json: km`
	AKKU  int16      `json: akku`
	STATE bool       `json:state`
}

type scootersHandlers struct {
	sync.Mutex
	store map[string]EScooter
}

func newScootersHandlers() *scootersHandlers {

	return &scootersHandlers{
		store: setUp(JSONFILE),
	}

}

func (h *scootersHandlers) scooters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return

	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}
}

func (h *scootersHandlers) post(w http.ResponseWriter, r *http.Request) {
	respBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	/*
		//TODO maybe not neccesary
		if r.Header.Get("Content-Type") != "application/json" {

			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Use only application/json"))
			return

		}
	*/
	var scooter EScooter
	err = json.Unmarshal(respBody, &scooter)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	s := h.store[scooter.ID]

	//TODO check dereference of pointer by post request without gps field
	if scooter.GPS == nil {
		scooter.GPS = s.GPS
	}
	if scooter.AKKU == 0 {
		scooter.AKKU = s.AKKU
	}
	if scooter.KM == 0 {
		scooter.KM = s.KM
	}

	h.Lock()
	h.store[scooter.ID] = scooter
	defer h.Unlock()

}

func (h *scootersHandlers) get(w http.ResponseWriter, r *http.Request) {

	scooters := make([]EScooter, len(h.store))

	h.Lock()
	i := 0
	for _, scooter := range h.store {
		scooters[i] = scooter
		i++
	}
	h.Unlock()

	jsonBytes, err := json.Marshal(scooters)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h *scootersHandlers) getScooter(w http.ResponseWriter, r *http.Request) {

	url := strings.Split(r.URL.String(), "/")
	if len(url) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()

	scooter, ok := h.store[url[2]]

	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonBytes, err := json.Marshal(scooter)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}

func setUp(j []byte) map[string]EScooter {

	var fJSON []EScooter

	e := json.Unmarshal(j, &fJSON)

	if e != nil {

		fmt.Println(e.Error())

	}

	scooters := make(map[string]EScooter, len(fJSON))

	for _, v := range fJSON {

		scooters[v.ID] = v

	}

	return scooters

}

func main() {

	sMap := newScootersHandlers()

	http.HandleFunc("/api", sMap.scooters)
	http.HandleFunc("/api/", sMap.getScooter)

	log.Fatal(http.ListenAndServe(":10000", nil))

}
