package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

func PointHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-type", "application/json")
	switch r.Method {
	case "GET":
		ts := r.URL.Path[len("/api/v1/point/"):]
		if len(ts) < 1 {
			http.Error(w, "Invalid parameter range", http.StatusBadRequest)
			return
		}
		log.Println(ts)

		BPDB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("BPBucket"))
			v := b.Get([]byte(ts))
			if err != nil {
				http.Error(w, "fetching bp", http.StatusInternalServerError)
				return err
			}
			if v == nil {
				http.Error(w, "fetching bp", http.StatusNotFound)
				return nil
			}
			fmt.Fprintf(w, string(v))
			return nil
		})
		break
	case "POST":
		payload, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err != nil {
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}

		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if payload == nil {
			http.Error(w, "No payload", http.StatusNotFound)
			return
		}
		bp := BloodPressurePoint{}
		err = json.Unmarshal(payload, &bp)
		ts := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
		bp.Timestamp = ts

		BPDB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("BPBucket"))
			id, err := b.NextSequence()
			if err != nil {
				http.Error(w, "fetching new id", http.StatusInternalServerError)
				return err
			}
			uid := int(id)
			bp.Id = uid
			bb, err := bp.Bytes()
			if err != nil {
				http.Error(w, "fetching new id", http.StatusInternalServerError)
				return err
			}

			err = b.Put([]byte(ts), bb)
			if err != nil {
				http.Error(w, "creating bp", http.StatusInternalServerError)
				return err
			}
			fmt.Fprintf(w, strconv.Itoa(uid))
			return nil
		})

		break
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func RangeHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		log.Println(r.URL.Path)
		// start_ts, end_ts
		start_ts := r.URL.Query().Get("start_ts")
		end_ts := r.URL.Query().Get("end_ts")
		if len(start_ts) < 1 || len(end_ts) < 1 {
			http.Error(w, "No time range", http.StatusBadRequest)
			return
		}
		bp := BloodPressureRange{}

		BPDB.View(func(tx *bolt.Tx) error {
			c := tx.Bucket([]byte("BPPoints")).Cursor()

			start := []byte(start_ts)
			end := []byte(end_ts)
			nn := 0
			for k, v := c.Seek(start); k != nil && bytes.Compare(k, end) <= 0; k, v = c.Next() {
				pp := BloodPressurePoint{}
				err = json.Unmarshal(v, &pp)
				if err != nil {
					http.Error(w, "unmarshalling pp", http.StatusInternalServerError)
					return err
				}

				bp.Points = append(bp.Points, pp)
				fmt.Printf("%s: %s\n", k, v)
				nn++
			}

			return nil
		})
		bs, err := bp.Bytes()
		if err != nil {
			http.Error(w, "bp unmarshall", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(bs))
		return
		break
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
