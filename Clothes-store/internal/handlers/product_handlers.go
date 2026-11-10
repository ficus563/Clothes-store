package handlers

import (
	"encoding/json"
	"net/http"
	"clothes-shop/internal/models"
	"clothes-shop/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

// продукты
// обработчик создания товара
func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Неверный JSON", http.StatusBadRequest)
		return
	}
	p.ID = uuid.New()
	if err := h.repo.CreateProduct(r.Context(), p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// обработчик получения всех товаров
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(products)
}

// обработчик удаления товара
func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := uuid.Parse(idStr)
	if err := h.repo.DeleteProduct(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Удалено"}`))
}

// роли
// обработчик создания роли
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	json.NewDecoder(r.Body).Decode(&role)
	role.ID = uuid.New()
	h.repo.CreateRole(r.Context(), role)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

// корзина
// обработчик добавления в корзину
func (h *Handler) CreateCart(w http.ResponseWriter, r *http.Request) {
	var c models.Cart
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = uuid.New()
	h.repo.CreateCart(r.Context(), c)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

// заказы
// обработчик создания заказа
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	json.NewDecoder(r.Body).Decode(&o)
	o.ID = uuid.New()
	h.repo.CreateOrder(r.Context(), o)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

// обработчик получения всех заказов
func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, _ := h.repo.GetAllOrders(r.Context())
	json.NewEncoder(w).Encode(orders)
}