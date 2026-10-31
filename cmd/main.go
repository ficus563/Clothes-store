package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"

	"clothes-shop/internal/handlers"
	"clothes-shop/internal/repository"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, используем системные переменные")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./clothes_shop.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("Ошибка подключения к БД", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("База данных не отвечает", "error", err)
		os.Exit(1)
	}
	slog.Info("Успешное подключение к SQLite")

	sqlBytes, err := os.ReadFile("migrations/ScriptBd.sql")
	if err != nil {
		slog.Error("Ошибка чтения SQL-скрипта", "error", err)
		os.Exit(1)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		slog.Error("Ошибка выполнения SQL-скрипта", "error", err)
		os.Exit(1)
	}
	slog.Info("Таблицы успешно созданы/проверены")

	productRepo := repository.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	r := chi.NewRouter()

	r.Post("/product", productHandler.Create)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Сервер работает..."))
	})

	slog.Info("Сервер запускается", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Ошибка сервера", "error", err)
		os.Exit(1)
	}
}
