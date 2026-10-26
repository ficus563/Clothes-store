package handlers

import (
	"endcoding/json"
	"net/http"
	"clothes-shop/internal/models"
	"clothes-shop/internal/repository"
	"github.com/google/uuid"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler {repo: repo}
}

func (h *ProductHandler) Create(w http.ResponseWriter,  r *http.Requset) {
	var p models.Product
	//чтение json
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Неверный формат Json", http.StatusBadRequest)
		return
	}
	p.ID = uuid.New()
	err = h.repo.Create(r.Context(), p)
	if err != nil {
		http.Error(w, "Ошибка при сохранении в базу", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}
//Todo
