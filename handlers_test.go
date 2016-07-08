package main

import (
	"bytes"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

func createFakeBP() BloodPressurePoint {
	bp := BloodPressurePoint{}
	ts := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	bp.Timestamp = ts
	bp.Description = "Linus"
	bp.Systolic = (rand.Float32() * 12) + 12.0
	bp.Diastolic = (rand.Float32() * 8) + 8.0
	bp.Tags = []string{"bbq"}
	return bp
}

func TestMain(m *testing.M) {
	BPDB, err = bolt.Open("test_bp.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	for a := 0; a < 20; a++ {
		bp := createFakeBP()
		log.Println(bp.Id)
	}
	r := m.Run()
	os.Exit(r)
}

func TestGETPointHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/point/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PointHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: %v (correct: %v)",
			status, http.StatusOK)
	}

	if len(rr.Body.String()) < 5 {
		t.Errorf("handler returned unexpected body: got %v ",
			rr.Body.String())
	}
}

func TestPOSTPointHandler(t *testing.T) {

	bp := createFakeBP()
	bb, err := bp.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/api/v1/point", bytes.NewBuffer(bb))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PointHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Wrong status code: %v (correct: %v)",
			status, http.StatusOK)
	}
	if rr.Body.String() != "11" {
		t.Errorf("handler returned unexpected body: found props: %v ",
			rr.Body.String())
	}
}
