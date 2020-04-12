package model

import (
	"encoding/json"
	"strings"
	"time"
)

type DumpDb struct {
	LastEditID int           `json:"last_edit_id"`
	Include    DumpIncludes  `json:"include"`
	Data       DumpGamesData `json:"data"`
}

// Platforms

type DumpIncludes struct {
	Platform DumpPlatformsData `json:"platform"`
}

type DumpPlatformsData struct {
	Data map[string]DumpPlatform `json:"data"`
}

type DumpPlatform struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
}

// Games

type DumpGamesData struct {
	Games []DumpGame `json:"games"`
}

type DumpGame struct {
	ID           int              `json:"id"`
	GameTitle    string           `json:"game_title"`
	PlatformID   int              `json:"platform"`
	ReleaseDate  *DumpGamesDbDate `json:"release_date"`
	Overview     *string          `json:"overview"`
	Youtube      *string          `json:"youtube"`
	Players      *int             `json:"players"`
	Coop         *string          `json:"coop"`
	Rating       *string          `json:"rating"`
	DeveloperIDS *[]int           `json:"developers"`
	GenreIDS     *[]int           `json:"genres"`
	PublisherIDS *[]int           `json:"publishers"`
	Alternatives *[]string        `json:"alternatives"`
	Uids         *[]DumpUIDType   `json:"uids"`
}

type DumpUIDType struct {
	UID                 string `json:"uid"`
	GamesUidsPatternsID int    `json:"games_uids_patterns_id"`
}

type DumpGamesDbDate struct {
	time.Time
}

func (sd *DumpGamesDbDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}

	sd.Time = newTime
	return nil
}

func (sd *DumpGamesDbDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(sd.Format("2006-01-02"))
}
