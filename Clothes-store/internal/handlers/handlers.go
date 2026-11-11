package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"clothes-shop/internal/models"
	"clothes-shop/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := r.Header.Get("X-User-Role")
			allowed := false
			for _, role := range roles {
				if userRole == role {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Доступ запрещен", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// --- АВТОРИЗАЦИЯ ---
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	// Очищаем от случайных пробелов
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	u, roleName, err := h.repo.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		fmt.Println("❌ Ошибка БД при входе:", err)
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	if u == nil {
		fmt.Println("❌ Пользователь не найден в БД:", req.Email)
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		fmt.Println("❌ Неверный пароль для:", req.Email)
		http.Error(w, "Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	fmt.Println("✅ Успешный вход:", req.Email, "| Роль:", roleName)
	json.NewEncoder(w).Encode(models.LoginResponse{ID: u.ID, Name: u.Name, RoleName: roleName})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)

	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hash)
	u.ID = uuid.New()
	u.RegistrationDate = time.Now()

	if u.RoleID == uuid.Nil {
		u.RoleID = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	}

	if err := h.repo.CreateUser(r.Context(), u); err != nil {
		fmt.Println("❌ Ошибка при регистрации:", err)
		http.Error(w, "Email уже занят", http.StatusBadRequest)
		return
	}
	u.Password = ""
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

// --- ТОВАРЫ ---
func (h *Handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}

	pType := r.URL.Query().Get("type")
	gender := r.URL.Query().Get("gender")
	brand := r.URL.Query().Get("brand")
	sortPrice := r.URL.Query().Get("sort")

	ps, err := h.repo.GetProducts(r.Context(), pType, gender, brand, limit, (page-1)*limit, sortPrice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ps)
}

func (h *Handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	p, _ := h.repo.GetProductByID(r.Context(), id)
	if p == nil {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	json.NewDecoder(r.Body).Decode(&p)
	p.ID = uuid.New()
	h.repo.CreateProduct(r.Context(), p)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var p models.Product
	json.NewDecoder(r.Body).Decode(&p)
	p.ID = id
	h.repo.UpdateProduct(r.Context(), p)
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.repo.DeleteProduct(r.Context(), id)
	w.WriteHeader(http.StatusOK)
}

// --- ФОТО ---
func (h *Handler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Нет файла", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, _ := io.ReadAll(file)
	hash := sha256.Sum256(fileBytes)
	hashStr := hex.EncodeToString(hash[:]) + ".jpg"

	if existing, _ := h.repo.GetPhotoByFilename(r.Context(), hashStr); existing != nil {
		json.NewEncoder(w).Encode(map[string]string{"message": "Файл уже есть", "id": existing.ID.String(), "filename": hashStr})
		return
	}

	os.MkdirAll("./uploads", os.ModePerm)
	out, _ := os.Create(filepath.Join("./uploads", hashStr))
	out.Write(fileBytes)
	out.Close()

	productID, _ := uuid.Parse(r.FormValue("product_id"))
	photo := models.Photo{ID: uuid.New(), ProductID: productID, Filename: hashStr}
	h.repo.CreatePhoto(r.Context(), photo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

func (h *Handler) GetPhoto(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	p, _ := h.repo.GetPhotoByID(r.Context(), id)
	if p == nil {
		http.Error(w, "Нет фото", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filepath.Join("./uploads", p.Filename))
}

func (h *Handler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.repo.DeletePhoto(r.Context(), id)
	w.WriteHeader(http.StatusOK)
}

// --- РОЛИ ---
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	json.NewDecoder(r.Body).Decode(&role)
	role.ID = uuid.New()
	h.repo.CreateRole(r.Context(), role)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (h *Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	rs, _ := h.repo.GetRoles(r.Context())
	json.NewEncoder(w).Encode(rs)
}

func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var role models.Role
	json.NewDecoder(r.Body).Decode(&role)
	role.ID = id
	h.repo.UpdateRole(r.Context(), role)
	json.NewEncoder(w).Encode(role)
}

func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.repo.DeleteRole(r.Context(), id)
	w.WriteHeader(http.StatusOK)
}

// --- ПОЛЬЗОВАТЕЛИ ---
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	us, _ := h.repo.GetUsers(r.Context(), limit, (page-1)*limit)
	json.NewEncoder(w).Encode(us)
}

func (h *Handler) GetUsersByRole(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	us, _ := h.repo.GetUsersByRole(r.Context(), id)
	json.NewEncoder(w).Encode(us)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)
	u.ID = id
	h.repo.UpdateUser(r.Context(), u)
	json.NewEncoder(w).Encode(u)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.repo.DeleteUser(r.Context(), id)
	w.WriteHeader(http.StatusOK)
}

// --- КОРЗИНА ---
func (h *Handler) CreateCart(w http.ResponseWriter, r *http.Request) {
	var c models.Cart
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = uuid.New()
	h.repo.CreateCart(r.Context(), c)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := uuid.Parse(chi.URLParam(r, "user_id"))
	cart, _ := h.repo.GetCartByUserID(r.Context(), userID)
	json.NewEncoder(w).Encode(cart)
}

func (h *Handler) UpdateCart(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var c models.Cart
	json.NewDecoder(r.Body).Decode(&c)
	c.ID = id
	h.repo.UpdateCart(r.Context(), c)
	json.NewEncoder(w).Encode(c)
}

func (h *Handler) DeleteCart(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	h.repo.DeleteCart(r.Context(), id)
	w.WriteHeader(http.StatusOK)
}

// --- ЗАКАЗЫ ---
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	json.NewDecoder(r.Body).Decode(&o)
	o.ID = uuid.New()
	o.DateTime = time.Now()
	if o.Status == "" {
		o.Status = models.StatusAccepted
	}
	h.repo.CreateOrder(r.Context(), o)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	orders, _ := h.repo.GetOrders(r.Context(), limit, (page-1)*limit)
	json.NewEncoder(w).Encode(orders)
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	o, _ := h.repo.GetOrderByID(r.Context(), id)
	if o == nil {
		http.Error(w, "Не найдено", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(o)
}

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	var o models.Order
	json.NewDecoder(r.Body).Decode(&o)
	o.ID = id
	h.repo.UpdateOrder(r.Context(), o)
	json.NewEncoder(w).Encode(o)
}

// --- ЭЛЕМЕНТЫ ЗАКАЗА ---
func (h *Handler) CreateOrderItem(w http.ResponseWriter, r *http.Request) {
	var item models.OrderItem
	json.NewDecoder(r.Body).Decode(&item)
	item.ID = uuid.New()
	h.repo.CreateOrderItem(r.Context(), item)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (h *Handler) GetOrderItems(w http.ResponseWriter, r *http.Request) {
	orderID, _ := uuid.Parse(chi.URLParam(r, "order_id"))
	items, _ := h.repo.GetOrderItems(r.Context(), orderID)
	json.NewEncoder(w).Encode(items)
}
