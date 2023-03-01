package main

import (
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type HTTPDriver struct {
	started bool
	store   *Store
}

func NewHTTPDriver(s *Store) *HTTPDriver {
	return &HTTPDriver{
		store: s,
	}
}

func (d *HTTPDriver) Run() {
	if d.started {
		return
	}

	router := httprouter.New()
	router.GET("/get/:key", d.GetHandler)
	router.GET("/has/:key", d.HasHandler)
	router.PUT("/put/:key", d.SetHandler)

	d.started = true
	log.Fatal(http.ListenAndServe(":7001", router))
}

func (d *HTTPDriver) GetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	w.Write(d.store.get(key))
}

func (d *HTTPDriver) HasHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	if e := d.store.has(key); e {
		w.WriteHeader(http.StatusOK)
	}
	w.WriteHeader(http.StatusNotFound)
}

func (d *HTTPDriver) SetHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	value, _ := io.ReadAll(r.Body)
	d.store.set(key, value)
	w.WriteHeader(http.StatusOK)
}
