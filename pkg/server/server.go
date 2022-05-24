package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/J-Swift/GamesDbMirror-go/pkg/model"
	"github.com/J-Swift/GamesDbMirror-go/pkg/repo"
	"github.com/gorilla/mux"
)

const (
	parsedDumpFilePath = "_clean.json"
)

var db *repo.Repo
var baseUrls map[string]string

type GamesResult struct {
	Games         []model.Game
	ImageBaseUrls map[string]string
}

func prepare(dataDir string) {
	filepath := dataDir + "/" + parsedDumpFilePath
	cleanDb := parse(filepath)
	db = repo.New(cleanDb.Games)
	baseUrls = cleanDb.ImageBaseUrls
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parse(dataFilepath string) model.CleanDB {
	fmt.Println("Parsing games")

	cachedDb, err := ioutil.ReadFile(dataFilepath)
	check(err)

	var result model.CleanDB
	json.Unmarshal(cachedDb, &result)

	fmt.Printf("  -> parsed [%d] games\n", len(result.Games))
	fmt.Println("  -> done")
	return result
}

func getIds(str string) []int {
	var result []int
	for _, id := range strings.Split(str, ",") {
		intID, _ := strconv.Atoi(id)
		result = append(result, intID)
	}
	return result
}

func getInsensitive(params url.Values, key string) string {
	for k, vs := range params {
		if strings.EqualFold(k, key) {
			if len(vs) == 0 {
				return ""
			}
			return vs[0]
		}
	}

	return ""
}

func handleFindByName(maxResults int) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		name := getInsensitive(req.URL.Query(), "name")

		games := db.FindGamesByTitle(name, maxResults)

		json.NewEncoder(w).Encode(GamesResult{Games: games, ImageBaseUrls: baseUrls})
	})
}

func handleFindByID(maxResults int) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		qIds := getInsensitive(req.URL.Query(), "ids")

		ids := getIds(qIds)
		games := db.FindGamesByID(ids, maxResults)

		json.NewEncoder(w).Encode(GamesResult{Games: games, ImageBaseUrls: baseUrls})
	})
}

func jsonContentMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, req)
	})
}

func Run(dataDir string, port string, maxResultsPerRequest int) {
	prepare(dataDir)

	fmt.Println("Initializing router")

	r := mux.NewRouter()
	r.Use(jsonContentMiddleware)

	gamesRouter := r.PathPrefix("/Games").Subrouter()
	gamesRouter.HandleFunc("/ByName", handleFindByName(maxResultsPerRequest)).Methods("GET")
	gamesRouter.HandleFunc("/ByIds", handleFindByID(maxResultsPerRequest)).Methods("GET")

	fmt.Println("  -> done")

	fmt.Printf("Now listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
