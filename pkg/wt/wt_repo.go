package wt

import (
	"github.com/4d3v/gowtbot/pkg/config"
	"github.com/4d3v/gowtbot/pkg/db"
)

type WTRepo struct {
	App           *config.AppConfig
	DBRepo        db.DBRepo
	Phones        map[string]bool
	BlockedPhones map[string]bool
}

func NewWTRepo(app *config.AppConfig,
	dbRepo db.DBRepo,
	phones map[string]bool,
	blockedPhones map[string]bool,
) *WTRepo {
	return &WTRepo{App: app,
		DBRepo:        dbRepo,
		Phones:        phones,
		BlockedPhones: blockedPhones,
	}
}

var wtRepo *WTRepo

func NewHandlers(repo *WTRepo) {
	wtRepo = repo
}
