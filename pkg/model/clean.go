package model

import (
	"strconv"
	"strings"
)

type CleanDB struct {
	Games         []Game
	ImageBaseUrls map[string]string
}

type Game struct {
	ID             int
	Title          string
	Platform       platform
	ReleaseDate    *DumpGamesDbDate
	Overview       NullString
	Youtube        NullString
	Players        NullInt
	IsCoop         NullBool
	Rating         NullString
	DeveloperIDS   []int
	GenreIDS       []int
	PublisherIDS   []int
	AlternateNames []string
	Uids           []uidType
	Images         []image
}

type platform struct {
	ID    int
	Name  string
	Alias string
}

type image struct {
	Id         int
	Type       string
	Side       NullString
	Filename   string
	Resolution NullString
}

type uidType struct {
	UID                 string
	GamesUidsPatternsID int
}

func NewGame(db *DumpDb, source *DumpGame) Game {
	platLookup := db.Include.Platform.ByGameId[strconv.Itoa(source.PlatformID)]
	imagesLookup, foundImages := db.Include.Images.ByGameId[strconv.Itoa(source.ID)]

	plat := platform{
		ID:    platLookup.ID,
		Name:  platLookup.Name,
		Alias: platLookup.Alias,
	}

	dids := []int{}
	if source.DeveloperIDS != nil {
		dids = *source.DeveloperIDS
	}
	gids := []int{}
	if source.GenreIDS != nil {
		gids = *source.GenreIDS
	}
	pids := []int{}
	if source.PublisherIDS != nil {
		pids = *source.PublisherIDS
	}
	alts := []string{}
	if source.Alternatives != nil {
		alts = *source.Alternatives
	}
	uids := []uidType{}
	if source.Uids != nil {
		for _, uid := range *source.Uids {
			uids = append(uids, uidType{
				UID:                 uid.UID,
				GamesUidsPatternsID: uid.GamesUidsPatternsID,
			})
		}
	}

	images := []image{}
	if foundImages {
		for _, img := range imagesLookup {
			var side NullString
			if img.Side != nil && *img.Side != "" {
				side = NullString{*img.Side, true}
			}
			var resolution NullString
			if img.Resolution != nil && *img.Resolution != "" {
				resolution = NullString{*img.Resolution, true}
			}
			images = append(images, image{
				Id:         img.Id,
				Type:       img.Type,
				Side:       side,
				Filename:   img.Filename,
				Resolution: resolution,
			})
		}
	}

	var ov NullString
	if source.Overview != nil && *source.Overview != "" {
		ov = NullString{*source.Overview, true}
	}
	var yt NullString
	if source.Youtube != nil && *source.Youtube != "" {
		yt = NullString{*source.Youtube, true}
	}
	var ps NullInt
	if source.Players != nil && *source.Players != 0 {
		ps = NullInt{int32(*source.Players), true}
	}
	var co NullBool
	if source.Coop != nil {
		co = NullBool{strings.ToLower(*source.Coop) == "yes", true}
	}
	var rt NullString
	if source.Rating != nil {
		rt = NullString{*source.Rating, true}
	}

	return Game{
		ID:             source.ID,
		Title:          source.GameTitle,
		Platform:       plat,
		ReleaseDate:    source.ReleaseDate,
		Overview:       ov,
		Youtube:        yt,
		Players:        ps,
		IsCoop:         co,
		Rating:         rt,
		DeveloperIDS:   dids,
		GenreIDS:       gids,
		PublisherIDS:   pids,
		AlternateNames: alts,
		Uids:           uids,
		Images:         images,
	}
}
