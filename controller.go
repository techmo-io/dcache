package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"log"
	"net/http"
)

type Controller struct {
	cluster *Cluster
}

func NewController(host Host) *Controller {
	controller := Controller{
		cluster: NewCluster(DefaultClusterConfig),
	}

	router := httprouter.New()
	router.GET("/join/:host/:port", controller.JoinHandler)
	go func() {
		log.Fatal(http.ListenAndServe(":"+host.ControllerPort, router))
	}()

	return &controller
}

// Join contacts the specified host, requesting to join its cluster
func (c *Controller) Join(host Host) error {
	r, e := http.Get(fmt.Sprint(host))
	if e != nil {
		return fmt.Errorf("join error: %w", e)
	}

	bytes, _ := io.ReadAll(r.Body)
	fmt.Printf("\nresponse:\n%s", bytes)
	if e := json.Unmarshal(bytes, c.cluster); e != nil {
		return fmt.Errorf("join error in unmarshal: %w", e)
	}

	return nil
}

func (d *Controller) JoinHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	host := Host{
		Name: ps.ByName("name"),
		Port: ps.ByName("port"),
	}
	d.cluster.AddNode(NewNode(host))
}
