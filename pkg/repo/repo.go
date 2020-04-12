package repo

import (
	"strings"

	"github.com/J-Swift/GamesDbMirror-go/pkg/model"
)

// Repo is an abstraction around games storage and retrieval
type Repo struct {
	games []model.Game
}

// New generates a new Repo
func New(games []model.Game) *Repo {
	return &Repo{
		games,
	}
}

// FindGamesByID queries the repo for the given game ids
func (repo Repo) FindGamesByID(ids []int, limit int) []model.Game {
	lookup := buildLookup(ids)
	return repo.findGames(limit, func(g model.Game) bool {
		return lookup[g.ID]
	})
}

// FindGamesByTitle queries the repo for games whos title partially matches the given title
func (repo *Repo) FindGamesByTitle(title string, limit int) []model.Game {
	title = strings.ToLower(title)
	return repo.findGames(limit, func(g model.Game) bool {
		return strings.Contains(strings.ToLower(g.Title), title)
	})
}

// Helpers

func buildLookup(elements []int) map[int]bool {
	result := make(map[int]bool)

	for _, el := range elements {
		result[el] = true
	}

	return result
}

func (repo Repo) findGames(limit int, isMatch func(model.Game) bool) []model.Game {
	if limit < 1 {
		limit = 9999999
	}

	result := []model.Game{}
	totalFound := 0

	for _, g := range repo.games {
		if isMatch(g) {
			result = append(result, g)
			totalFound++
			if totalFound == limit {
				break
			}
		}
	}

	return result
}
