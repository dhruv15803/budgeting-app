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

// parseDate parses a YYYY-MM-DD string into a time.Time pointer; returns nil on empty input.
func parseDate(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

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

	q := r.URL.Query()
	f := models.ExpenseFilter{
		Page:     1,
		PageSize: 20,
		SortBy:   q.Get("sort_by"),
	}

	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			f.Page = v
		}
	}
	if ps := q.Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			f.PageSize = v
		}
	}

	if s := q.Get("search"); s != "" {
		f.Search = &s
	}

	// Explicit date range takes precedence over month/year shorthands.
	dateFrom, err := parseDate(q.Get("date_from"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "date_from must be in YYYY-MM-DD format")
		return
	}
	dateTo, err := parseDate(q.Get("date_to"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "date_to must be in YYYY-MM-DD format")
		return
	}

	// month shorthand: ?month=YYYY-MM expands to first and last day of the month.
	if dateFrom == nil && dateTo == nil {
		if monthStr := q.Get("month"); monthStr != "" {
			m, err := time.Parse("2006-01", monthStr)
			if err != nil {
				_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
				return
			}
			first := time.Date(m.Year(), m.Month(), 1, 0, 0, 0, 0, time.UTC)
			last := first.AddDate(0, 1, -1)
			dateFrom = &first
			dateTo = &last
		} else if yearStr := q.Get("year"); yearStr != "" {
			// year shorthand: ?year=YYYY expands to Jan 1 – Dec 31.
			yr, err := strconv.Atoi(yearStr)
			if err != nil || yr < 1 {
				_ = writeJsonError(w, http.StatusBadRequest, "year must be a valid 4-digit year")
				return
			}
			first := time.Date(yr, time.January, 1, 0, 0, 0, 0, time.UTC)
			last := time.Date(yr, time.December, 31, 0, 0, 0, 0, time.UTC)
			dateFrom = &first
			dateTo = &last
		}
	}

	f.DateFrom = dateFrom
	f.DateTo = dateTo

	// category_id is repeatable: ?category_id=1&category_id=3
	for _, cidStr := range q["category_id"] {
		if cid, err := strconv.Atoi(cidStr); err == nil && cid > 0 {
			f.CategoryIDs = append(f.CategoryIDs, cid)
		}
	}

	if amtMin := q.Get("amount_min"); amtMin != "" {
		v, err := strconv.ParseFloat(amtMin, 64)
		if err != nil || v < 0 {
			_ = writeJsonError(w, http.StatusBadRequest, "amount_min must be a non-negative number")
			return
		}
		f.AmountMin = &v
	}
	if amtMax := q.Get("amount_max"); amtMax != "" {
		v, err := strconv.ParseFloat(amtMax, 64)
		if err != nil || v < 0 {
			_ = writeJsonError(w, http.StatusBadRequest, "amount_max must be a non-negative number")
			return
		}
		f.AmountMax = &v
	}

	expenses, total, err := h.services.Expenses.ListExpenses(claims.UserID, f)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	items := make([]expenseResponse, 0, len(expenses))
	for i := range expenses {
		items = append(items, toExpenseResponse(&expenses[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(f.PageSize)))
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
			Page:       f.Page,
			PageSize:   f.PageSize,
			TotalPages: totalPages,
		},
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
