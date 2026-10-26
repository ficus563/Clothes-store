package main 

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"modernc.org/sqlite"
	"clother-store/internal/handlers"
	"clothes-store/internal/repository"
)

func main( {
	var logger = slog.New(slog.NewJSONHandlers(os.Stdout, nil))
	slog.SetDefault(logger)
	if err := godotenv.Load(); err != nil {
		
	}
})

//TODO...
