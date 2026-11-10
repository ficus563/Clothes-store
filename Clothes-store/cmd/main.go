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
	// настраиваем логи
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// загружаем переменные из .env
	if err := godotenv.Load(); err != nil {
		slog.Warn("Файл .env не найден, используем стандартные настройки")
	}

	// получаем порт и путь к базе
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./clothes_shop.db"
	}

	// подключаем базу
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		slog.Error("Ошибка подключения к БД", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// читаем sql скрипт и создаем таблицы
	sqlBytes, err := os.ReadFile("migrations/ScriptBd.sql")
	if err != nil {
		slog.Error("Ошибка чтения SQL-скрипта", "error", err)
		os.Exit(1)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		slog.Error("Ошибка создания таблиц", "error", err)
		os.Exit(1)
	}

	// инициализируем репозиторий и хендлеры
	repo := repository.NewRepository(db)
	h := handlers.NewHandler(repo)

	// настраиваем роутер
	r := chi.NewRouter()

	// маршруты для товаров
	r.Post("/product", h.CreateProduct)
	r.Get("/product", h.GetProducts)
	r.Delete("/product/{id}", h.DeleteProduct)

	// маршруты для ролей, корзины и заказов
	r.Post("/role", h.CreateRole)
	r.Post("/cart", h.CreateCart)
	r.Post("/order", h.CreateOrder)
	r.Get("/order", h.GetOrders)

	// раздаем html файл
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// запускаем сервер
	slog.Info("Сервер запущен", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Ошибка сервера", "error", err)
	}
}