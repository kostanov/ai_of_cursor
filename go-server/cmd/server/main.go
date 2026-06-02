package main

import (
	"log"
	"net/http"

	"ai-of-cursor/go-server/internal/config"
	"ai-of-cursor/go-server/internal/database"
	"ai-of-cursor/go-server/internal/handler"
	"ai-of-cursor/go-server/internal/repository"
	"ai-of-cursor/go-server/internal/router"
	"ai-of-cursor/go-server/internal/service"
)

func main() {
	db, err := database.Open(config.DBPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err = database.InitSchema(db); err != nil {
		log.Fatalf("init db: %v", err)
	}

	users := repository.NewUserRepository(db)
	active := service.NewActiveUsers(config.MaxActiveUsers)
	h := handler.New(users, active)

	log.Printf("server started on %s", config.ServerAddr)
	if err = http.ListenAndServe(config.ServerAddr, router.New(h)); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
