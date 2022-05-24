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
	metaFilePath       = "_meta.json"
	parsedDumpFilePath = "_clean.json"
	dbDumpFilePath     = "_dump.json"

	dbDumpURL = "https://cdn.thegamesdb.net/json/database-latest.json"
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
		Version: 2,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readMeta(filepath string) dataFetchMeta {
	fmt.Println("Reading meta")

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Println("  -> no meta found")
		return newMeta()
	}

	res, err := ioutil.ReadFile(filepath)
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

func writeMeta(filepath string, meta dataFetchMeta) {
	fmt.Println("Writing meta")

	metaJSON, err := json.MarshalIndent(&meta, "", "  ")
	check(err)

	err = ioutil.WriteFile(filepath, metaJSON, 0644)
	check(err)

	fmt.Println("  -> done")
}

func downloadDbDump(filepath string, meta *dataFetchMeta) bool {
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
	err = ioutil.WriteFile(filepath, body, 0644)
	check(err)
	meta.SavedAt = time.Now().UTC()

	newEtag := resp.Header.Get("etag")
	meta.DbDumpEtag = newEtag
	fmt.Printf("  -> recorded etag [%s]\n", newEtag)
	fmt.Println("  -> done")

	return true
}

func parseDbDump(filepath string, meta *dataFetchMeta) model.CleanDB {
	fmt.Println("Parsing db")

	var data model.DumpDb

	s, err := ioutil.ReadFile(filepath)
	check(err)

	json.Unmarshal(s, &data)

	result := model.CleanDB{
		ImageBaseUrls: data.Include.Images.BaseUrls,
	}
	result.Games = make([]model.Game, len(data.Data.Games))
	for i, g := range data.Data.Games {
		result.Games[i] = model.NewGame(&data, &g)
	}

	lastEditID := data.LastEditID
	meta.DbDumpLastEditID = lastEditID
	fmt.Printf("  -> recorded DbDumpLastEditID [%d]\n", lastEditID)
	fmt.Println("  -> done")

	return result
}

func saveParsedDump(filepath string, parsed model.CleanDB) {
	fmt.Println("Saving parsed games")

	r, err := json.Marshal(parsed)
	check(err)

	err = ioutil.WriteFile(filepath, r, 0644)
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

func createOutDirIfNeeded(outDir string) {
	fmt.Println("Checking for out dir")

	info, err := os.Stat(outDir)
	if os.IsNotExist(err) {
		fmt.Println("  -> doesnt exist, creating")
		check(os.MkdirAll(outDir, 0755))
	} else if !info.IsDir() {
		check(fmt.Errorf("not a directory [%s]", outDir))
	}
	fmt.Println("  -> done")
}

func Run(outDir string) {
	metaPath := outDir + "/" + metaFilePath
	dbDumpPath := outDir + "/" + dbDumpFilePath
	parsedDumpPath := outDir + "/" + parsedDumpFilePath

	createOutDirIfNeeded(outDir)

	cachedMeta := readMeta(metaPath)

	currentMeta := cachedMeta
	currentMeta.RanAt = time.Now().UTC()

	if downloadDbDump(dbDumpPath, &currentMeta) {
		cleanDb := parseDbDump(dbDumpPath, &currentMeta)
		updateGames(&currentMeta, cleanDb.Games)
		saveParsedDump(parsedDumpPath, cleanDb)
	}

	writeMeta(metaPath, currentMeta)
}
