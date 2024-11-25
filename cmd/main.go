package main

import (
	"flag"
	"log/slog"
	"os"
	"todocli/internal/config"
	"todocli/internal/database"
	"todocli/internal/layout"
)

var (
	conf  config.Config
	debug bool
)

func main() {
	flag.StringVar(&conf.DBPath, "dbPath", os.Getenv("HOME")+"/.todocli.db", "Path to the sqlite database")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()

	if debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	c := conf.New()
	slog.Info("DB", "dbPath", c.DBPath)
	db := database.New(c.DBPath)
	if db == nil {
		slog.Error("Failed to initialize database")
		return
	}

	layout.UI(db)

	slog.Info("Closing database")
	err := db.Close()
	if err != nil {
		slog.Error("Error closing database", "error", err)
	}

}
