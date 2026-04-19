package handlers

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/services"
	"github.com/go-chi/chi/v5"
)

type createExpenseRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryID  int     `json:"category_id"`
	ExpenseDate string  `json:"expense_date"` // expected format: YYYY-MM-DD
}

type updateExpenseRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryID  int     `json:"category_id"`
	ExpenseDate string  `json:"expense_date"`
}

type expenseResponse struct {
	ID                 int        `json:"id"`
	Title              string     `json:"title"`
	Description        *string    `json:"description"`
	Amount             float64    `json:"amount"`
	UserID             int        `json:"user_id"`
	CategoryID         int        `json:"category_id"`
	RecurringExpenseID *int       `json:"recurring_expense_id"`
	ExpenseDate        string     `json:"expense_date"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

type expenseListData struct {
	Expenses   []expenseResponse `json:"expenses"`
	Total      int               `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

func toExpenseResponse(e *models.Expense) expenseResponse {
	return expenseResponse{
		ID:                 e.ID,
		Title:              e.Title,
		Description:        e.Description,
		Amount:             e.Amount,
		UserID:             e.UserID,
		CategoryID:         e.CategoryID,
		RecurringExpenseID: e.RecurringExpenseID,
		ExpenseDate:        e.ExpenseDate.Format("2006-01-02"),
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
}

func (h *Handler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "expense_date must be in YYYY-MM-DD format")
		return
	}

	expense, err := h.services.Expenses.CreateExpense(claims.UserID, req.CategoryID, req.Title, req.Description, req.Amount, expenseDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type resp struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    expenseResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusCreated, resp{
		Success: true,
		Message: "Expense created successfully",
		Data:    toExpenseResponse(expense),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	expenseID, err := strconv.Atoi(idStr)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	var req updateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	expenseDate, err := time.Parse("2006-01-02", req.ExpenseDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "expense_date must be in YYYY-MM-DD format")
		return
	}

	expense, err := h.services.Expenses.UpdateExpense(expenseID, claims.UserID, req.CategoryID, req.Title, req.Description, req.Amount, expenseDate)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, services.ErrForbidden) {
		_ = writeJsonError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type resp struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    expenseResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "Expense updated successfully",
		Data:    toExpenseResponse(expense),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	expenseID, err := strconv.Atoi(idStr)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid expense id")
		return
	}

	err = h.services.Expenses.DeleteExpense(expenseID, claims.UserID)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if errors.Is(err, services.ErrForbidden) {
		_ = writeJsonError(w, http.StatusForbidden, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJsonResponse(w, http.StatusOK, ApiResponse{
		Success: true,
		Message: "Expense deleted successfully",
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	page := 1
	pageSize := 20

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	expenses, total, err := h.services.Expenses.ListExpenses(claims.UserID, page, pageSize)
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]expenseResponse, 0, len(expenses))
	for i := range expenses {
		items = append(items, toExpenseResponse(&expenses[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	type resp struct {
		Success bool            `json:"success"`
		Message string          `json:"message"`
		Data    expenseListData `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "ok",
		Data: expenseListData{
			Expenses:   items,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
