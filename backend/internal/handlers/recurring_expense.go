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

type createRecurringExpenseRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryID  int     `json:"category_id"`
	StartDate   string  `json:"start_date"`         // YYYY-MM-DD
	EndDate     *string `json:"end_date,omitempty"` // YYYY-MM-DD or null
	Frequency   string  `json:"frequency"`          // daily|weekly|monthly|yearly
}

type updateRecurringExpenseRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Amount      float64 `json:"amount"`
	CategoryID  int     `json:"category_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	Frequency   string  `json:"frequency"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type recurringExpenseResponse struct {
	ID             int        `json:"id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description"`
	Amount         float64    `json:"amount"`
	UserID         int        `json:"user_id"`
	CategoryID     int        `json:"category_id"`
	StartDate      string     `json:"start_date"`
	EndDate        *string    `json:"end_date"`
	Frequency      string     `json:"frequency"`
	NextOccurrence string     `json:"next_occurrence"`
	IsActive       bool       `json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type recurringExpenseListData struct {
	RecurringExpenses []recurringExpenseResponse `json:"recurring_expenses"`
	Total             int                        `json:"total"`
	Page              int                        `json:"page"`
	PageSize          int                        `json:"page_size"`
	TotalPages        int                        `json:"total_pages"`
}

func toRecurringExpenseResponse(r *models.RecurringExpense) recurringExpenseResponse {
	var endDate *string
	if r.EndDate != nil {
		s := r.EndDate.Format("2006-01-02")
		endDate = &s
	}
	return recurringExpenseResponse{
		ID:             r.ID,
		Title:          r.Title,
		Description:    r.Description,
		Amount:         r.Amount,
		UserID:         r.UserID,
		CategoryID:     r.CategoryID,
		StartDate:      r.StartDate.Format("2006-01-02"),
		EndDate:        endDate,
		Frequency:      r.Frequency,
		NextOccurrence: r.NextOccurrence.Format("2006-01-02"),
		IsActive:       r.IsActive,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

// parseOptionalDate parses a YYYY-MM-DD pointer; nil or empty string returns nil.
func parseOptionalDate(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (h *Handler) CreateRecurringExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createRecurringExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "start_date must be in YYYY-MM-DD format")
		return
	}
	endDate, err := parseOptionalDate(req.EndDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "end_date must be in YYYY-MM-DD format")
		return
	}

	created, err := h.services.RecurringExpenses.CreateRecurringExpense(
		claims.UserID, req.CategoryID, req.Title, req.Description, req.Amount,
		startDate, endDate, req.Frequency,
	)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type resp struct {
		Success bool                     `json:"success"`
		Message string                   `json:"message"`
		Data    recurringExpenseResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusCreated, resp{
		Success: true,
		Message: "Recurring expense created successfully",
		Data:    toRecurringExpenseResponse(created),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) GetRecurringExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid recurring expense id")
		return
	}

	existing, err := h.services.RecurringExpenses.GetRecurringExpense(id, claims.UserID)
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

	type resp struct {
		Success bool                     `json:"success"`
		Message string                   `json:"message"`
		Data    recurringExpenseResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "ok",
		Data:    toRecurringExpenseResponse(existing),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) UpdateRecurringExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid recurring expense id")
		return
	}

	var req updateRecurringExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "start_date must be in YYYY-MM-DD format")
		return
	}
	endDate, err := parseOptionalDate(req.EndDate)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "end_date must be in YYYY-MM-DD format")
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	updated, err := h.services.RecurringExpenses.UpdateRecurringExpense(
		id, claims.UserID, req.CategoryID, req.Title, req.Description, req.Amount,
		startDate, endDate, req.Frequency, isActive,
	)
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
		Success bool                     `json:"success"`
		Message string                   `json:"message"`
		Data    recurringExpenseResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "Recurring expense updated successfully",
		Data:    toRecurringExpenseResponse(updated),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) DeleteRecurringExpense(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid recurring expense id")
		return
	}

	err = h.services.RecurringExpenses.DeleteRecurringExpense(id, claims.UserID)
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
		Message: "Recurring expense deleted successfully",
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) ListRecurringExpenses(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	q := r.URL.Query()
	f := models.RecurringExpenseFilter{
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
	for _, cidStr := range q["category_id"] {
		if cid, err := strconv.Atoi(cidStr); err == nil && cid > 0 {
			f.CategoryIDs = append(f.CategoryIDs, cid)
		}
	}
	if freq := q.Get("frequency"); freq != "" {
		f.Frequency = &freq
	}
	if active := q.Get("is_active"); active != "" {
		v, err := strconv.ParseBool(active)
		if err != nil {
			_ = writeJsonError(w, http.StatusBadRequest, "is_active must be true or false")
			return
		}
		f.IsActive = &v
	}

	items, total, err := h.services.RecurringExpenses.ListRecurringExpenses(claims.UserID, f)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	out := make([]recurringExpenseResponse, 0, len(items))
	for i := range items {
		out = append(out, toRecurringExpenseResponse(&items[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(f.PageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	type resp struct {
		Success bool                     `json:"success"`
		Message string                   `json:"message"`
		Data    recurringExpenseListData `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "ok",
		Data: recurringExpenseListData{
			RecurringExpenses: out,
			Total:             total,
			Page:              f.Page,
			PageSize:          f.PageSize,
			TotalPages:        totalPages,
		},
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
