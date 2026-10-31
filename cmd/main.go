package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite" // Драйвер SQLite (подключается через _)

	"clothes-shop/internal/handlers"   // Твои хендлеры (убедись, что путь совпадает с go.mod!)
	"clothes-shop/internal/repository" // Твои репозитории
)

func main() {
	//Настраиваем красивые JSON-логи (требование ТЗ)
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	//Загружаем настройки из .env
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

	//Подключаемся к базе данных
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

	//Читаем и применяем SQL-скрипт
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

	
	//Создаем репозиторий и передаем ему подключение к базе
	productRepo := repository.NewProductRepository(db)
	//Создаем хендлер и передаем ему репозиторий
	productHandler := handlers.NewProductHandler(productRepo)

	//Настройка роутера
	r := chi.NewRouter()

	//Регистрируем маршруты
	r.Post("/product", productHandler.Create)

	//Простой маршрут для проверки
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Сервер магазина одежды работает!"))
	})

	//Запуск сервера
	slog.Info("Сервер запускается", "port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Ошибка сервера", "error", err)
		os.Exit(1)
	}
}
