package main

//go:generate go-bindata-assetfs static/...

import (
	"log"
	"net/http"

	"go.uber.org/atomic"
)

const QUERY_PARAM = "query"
const SEARCH_URL = "https://lbry.tv/$/search?q="
const DEFAULT_BACKUP_QUERY = "LBRY HOT DANCE"

var previousQuery = atomic.NewString(DEFAULT_BACKUP_QUERY)

func main() {
	port := "8080"
	log.Println("Server started at localhost:" + port)
	http.Handle("/", http.FileServer(assetFS()))
	http.HandleFunc("/search", Handler)
	http.ListenAndServe(":"+port, nil)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()[QUERY_PARAM]
	if !ok {
		log.Println("No query", keys)
		return
	}

	searchQuery := keys[0]
	InsertAndRedirect(w, r, searchQuery)

	// if previousQuery != "" {
	// 	// We already have the previousQuery in memory, just return that
	// }

	// // First request since the server started, we need to go into the DB
	// var previousQueryInDb = Query{}
	// err := GetLatestValue(&previousQueryInDb)
	// if err != nil {
	// 	log.Println(err)
	// 	InsertAndRedirect(w, r, searchQuery, DEFAULT_BACKUP_QUERY)
	// 	return
	// }

	// InsertAndRedirect(w, r, searchQuery, previousQueryInDb.Value)
}

func InsertAndRedirect(w http.ResponseWriter, r *http.Request, currentQuery string) {
	s := previousQuery.Load()
	if s == "" {
		s = DEFAULT_BACKUP_QUERY
	}

	redirectUrl := SEARCH_URL + s

	previousQuery.Store(currentQuery)

	// Insert(currentQuery)

	http.Redirect(w, r, redirectUrl, 302)
}

//
//
//
// Database stuff
//
//
//

// type Query struct {
// 	ID        int    `json:"id"`
// 	Value     string `json:"value"`
// 	Timestamp string `json:",string"`
// }

// func dbConn() (db *sql.DB) {
// 	dbDriver := "mysql"
// 	dbUser := "root"
// 	dbPass := "password"
// 	dbName := "searchredirect"
// 	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	return db
// }

// func GetLatestValue(result *Query) error {
// 	db := dbConn()

// 	rows, err := db.Query("SELECT * FROM youtube ORDER BY timestamp DESC LIMIT 1")
// 	if err != nil {
// 		return err
// 	}

// 	// I know this will only be one row, I probably don't have to do this loop thing over it?
// 	for rows.Next() {
// 		var id int
// 		var value string
// 		var timestamp string

// 		err = rows.Scan(&id, &value, &timestamp)
// 		if err != nil {
// 			panic(err.Error())
// 		}

// 		result.Value = value
// 		result.Timestamp = timestamp
// 		result.ID = id
// 	}

// 	defer db.Close()
// 	return nil
// }

// func Insert(value string) error {
// 	db := dbConn()

// 	insForm, err := db.Prepare("INSERT INTO youtube(value) VALUES(?)")
// 	if err != nil {
// 		return err
// 	}

// 	_, err = insForm.Exec(value)
// 	if err != nil {
// 		log.Println("Error reading from db: ", err)
// 	}

// 	defer db.Close()
// 	return nil
// }
