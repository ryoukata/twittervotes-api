package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/tylerstillwater/graceful"
	"gopkg.in/mgo.v2"
)

func main() {
	var (
		addr = flag.String("addr", ":8080", "Address to Endpoint.")
		// mongo = flag.String("mongo", "twitter-votes-mongodb", "Address to MongoDB")
		mongo = flag.String("mongo", "localhost", "Address to MongoDB")
	)
	flag.Parse()
	log.Println("Connect to MongoDB...", *mongo)
	mongoInfo := &mgo.DialInfo{
		Addrs:    []string{"localhost:27017"},
		Timeout:  20 * time.Second,
		Database: "ballots",
		Username: "mongo",
		Password: "mongo",
		Source:   "ballots",
	}
	db, err := mgo.DialWithInfo(mongoInfo)
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB.: ", err)
	}
	// db, err := mgo.Dial(*mongo)
	// if err != nil {
	// 	log.Fatalln("Failed to connect to MongoDB.: ", err)
	// }
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	log.Println("Start Web Server.: ", *addr)
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("Stop Web Server...")
}

// check API Key.
func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "Invalid API Key.")
			return
		}
		fn(w, r)
	}
}

func isValidAPIKey(key string) bool {
	return key == "abc123"
}

// manage DB session.
func withData(d *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		f(w, r)
	}
}

// vars set up and clean up.
func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

// CORS settings
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
