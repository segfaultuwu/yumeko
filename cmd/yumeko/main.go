package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segfaultuwu/yumeko/internal/bot"
	"github.com/segfaultuwu/yumeko/internal/config"
	"github.com/segfaultuwu/yumeko/internal/database"
)

func main() {
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatal("config:", err)
	}

	db, err := database.Open(cfg.Database.Path)
	if err != nil {
		log.Fatal("database:", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal("migrate:", err)
	}

	yumeko, err := bot.New(cfg, db)
	if err != nil {
		log.Fatal("bot:", err)
	}

	if err := yumeko.Start(); err != nil {
		log.Fatal("start:", err)
	}

	log.Println("Yumeko is running. Press CTRL+C to stop.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	if err := yumeko.Stop(); err != nil {
		log.Println("stop:", err)
	}
}
