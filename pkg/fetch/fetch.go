package fetch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/J-Swift/GamesDbMirror-go/pkg/model"
)

const (
	metaFilePath = "out/_meta.json"

	dbDumpFilePath = "out/_dump.json"
	dbDumpURL      = "https://cdn.thegamesdb.net/json/database-latest.json"

	parsedDumpFilePath = "out/_clean.json"
)

type dataFetchMeta struct {
	Version          int
	RanAt            time.Time
	SavedAt          time.Time
	DbDumpEtag       string
	DbDumpLastEditID int
	GamesLastEditID  int
}

func newMeta() dataFetchMeta {
	return dataFetchMeta{
		Version: 1,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readMeta() dataFetchMeta {
	fmt.Println("Reading meta")

	if _, err := os.Stat(metaFilePath); os.IsNotExist(err) {
		fmt.Println("  -> no meta found")
		return newMeta()
	}

	res, err := ioutil.ReadFile(metaFilePath)
	check(err)

	var meta dataFetchMeta
	check(json.Unmarshal(res, &meta))

	if currentMeta := newMeta(); meta.Version != currentMeta.Version {
		fmt.Println("  -> stored meta is from older version")
		return currentMeta
	}

	fmt.Printf("  -> %+v\n", meta)
	fmt.Println("  -> done")
	return meta
}

func writeMeta(meta dataFetchMeta) {
	fmt.Println("Writing meta")

	metaJSON, err := json.MarshalIndent(&meta, "", "  ")
	check(err)

	err = ioutil.WriteFile(metaFilePath, metaJSON, 0644)
	check(err)

	fmt.Println("  -> done")
}

func downloadDbDump(meta *dataFetchMeta) bool {
	fmt.Println("Downloading db")

	req, err := http.NewRequest("GET", dbDumpURL, nil)
	check(err)
	if meta.DbDumpEtag != "" {
		req.Header.Add("If-None-Match", meta.DbDumpEtag)
	}

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	check(err)
	if resp.StatusCode == 304 {
		fmt.Println("  -> dump not modified since last download")
		return false
	}

	fmt.Println("  -> downloading content")
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	fmt.Println("  -> writing to disk")
	err = ioutil.WriteFile(dbDumpFilePath, body, 0644)
	check(err)
	meta.SavedAt = time.Now().UTC()

	newEtag := resp.Header.Get("etag")
	meta.DbDumpEtag = newEtag
	fmt.Printf("  -> recorded etag [%s]\n", newEtag)
	fmt.Println("  -> done")

	return true
}

func parseDbDump(meta *dataFetchMeta) []model.Game {
	fmt.Println("Parsing db")

	var data model.DumpDb

	s, err := ioutil.ReadFile(dbDumpFilePath)
	check(err)

	json.Unmarshal(s, &data)

	result := make([]model.Game, len(data.Data.Games))
	for i, g := range data.Data.Games {
		result[i] = model.NewGame(&data, &g)
	}

	lastEditID := data.LastEditID
	meta.DbDumpLastEditID = lastEditID
	fmt.Printf("  -> recorded DbDumpLastEditID [%d]\n", lastEditID)
	fmt.Println("  -> done")

	return result
}

func saveParsedDump(parsed []model.Game) {
	fmt.Println("Saving parsed games")

	r, err := json.Marshal(parsed)
	check(err)

	err = ioutil.WriteFile(parsedDumpFilePath, r, 0644)
	check(err)

	fmt.Println("  -> done")
}

func updateGames(meta *dataFetchMeta, games []model.Game) {
	fmt.Println("Updating games")

	lastEditID := meta.DbDumpLastEditID
	meta.GamesLastEditID = lastEditID
	fmt.Printf("  -> recorded GamesLastEditID [%d]\n", lastEditID)
	fmt.Println("  -> not implemented yet")
}

func Run() {
	cachedMeta := readMeta()

	currentMeta := cachedMeta
	currentMeta.RanAt = time.Now().UTC()

	if downloadDbDump(&currentMeta) {
		games := parseDbDump(&currentMeta)
		updateGames(&currentMeta, games)
		saveParsedDump(games)
	}

	writeMeta(currentMeta)
}
