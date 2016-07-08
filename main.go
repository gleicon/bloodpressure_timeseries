package main

import (
	"log"
	"net/http"

	"github.com/boltdb/bolt"
)

var BPDB *bolt.DB
var err error

func main() {
	BPDB, err = bolt.Open("bloodpoints.db", 0600, nil)
	if err != nil {
		log.Println(err)
		return
	}

	http.HandleFunc("/api/v1/point", PointHandler)
	http.HandleFunc("/api/v1/point/", PointHandler)
	http.HandleFunc("/api/v1/range", RangeHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
