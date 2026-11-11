package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"

	"clothes-shop/internal/handlers"
	"clothes-shop/internal/repository"
)

func main() {
	// логи
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	// база
	db, err := sql.Open("sqlite", "./clothes_shop.db")
	if err != nil {
		slog.Error("Ошибка БД", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// таблицы
	sqlBytes, _ := os.ReadFile("migrations/ScriptBd.sql")
	db.Exec(string(sqlBytes))

	// слои
	repo := repository.NewRepository(db)
	h := handlers.NewHandler(repo)

	// роутер
	r := chi.NewRouter()

	// --- API МАРШРУТЫ ---
	r.Post("/register", h.Register)
	r.Get("/product", h.GetProducts)
	r.Get("/product/{id}", h.GetProductByID)
	r.Get("/cart/{user_id}", h.GetCart)
	r.Post("/photo", h.UploadPhoto)
	r.Get("/order", h.GetOrders)

	// раздача картинок
	fs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", fs))

	// группы с ролями
	r.Group(func(r chi.Router) {
		r.Use(handlers.RoleMiddleware("Администратор"))
		r.Post("/product", h.CreateProduct)
		r.Put("/product/{id}", h.UpdateProduct)
		r.Delete("/product/{id}", h.DeleteProduct)
		r.Post("/role", h.CreateRole)
	})

	r.Group(func(r chi.Router) {
		r.Use(handlers.RoleMiddleware("Клиент"))
		r.Post("/cart", h.CreateCart)
		r.Delete("/cart/{id}", h.DeleteCart)
		r.Post("/order", h.CreateOrder)
	})

	// --- ГЛАВНАЯ СТРАНИЦА (HTML) ---
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// проверяем, есть ли файл в корне
		if _, err := os.Stat("index.html"); os.IsNotExist(err) {
			// если файла нет, выводим текст
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write([]byte("<h1>Сервер работает!</h1><p>Но файл index.html не найден в корне проекта.</p>"))
			return
		}
		// если файл есть, отдаем его
		http.ServeFile(w, r, "index.html")
	})

	// запуск
	slog.Info("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("Сервер упал", "error", err)
	}
}
