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

// parseMonth parses a "YYYY-MM" string and returns the first day of that month in UTC.
func parseMonth(s string) (time.Time, error) {
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

// ── request types ────────────────────────────────────────────────────────────

type createBudgetRequest struct {
	BudgetMonth string  `json:"budget_month"` // YYYY-MM
	TotalAmount float64 `json:"total_amount"`
}

type updateBudgetRequest struct {
	TotalAmount float64 `json:"total_amount"`
}

type upsertCategoryBudgetRequest struct {
	CategoryID      int     `json:"category_id"`
	AllocatedAmount float64 `json:"allocated_amount"`
}

type bulkSetCategoriesRequest struct {
	Categories []struct {
		CategoryID      int     `json:"category_id"`
		AllocatedAmount float64 `json:"allocated_amount"`
	} `json:"categories"`
}

// ── response types ───────────────────────────────────────────────────────────

type monthlyBudgetResponse struct {
	ID          int        `json:"id"`
	BudgetMonth string     `json:"budget_month"`
	TotalAmount float64    `json:"total_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type categoryBudgetResponse struct {
	ID              int     `json:"id"`
	CategoryID      int     `json:"category_id"`
	CategoryName    string  `json:"category_name,omitempty"`
	AllocatedAmount float64 `json:"allocated_amount"`
}

type categorySpendingResponse struct {
	ID              int     `json:"id"`
	CategoryID      int     `json:"category_id"`
	CategoryName    string  `json:"category_name"`
	AllocatedAmount float64 `json:"allocated_amount"`
	SpentAmount     float64 `json:"spent_amount"`
	Remaining       float64 `json:"remaining"`
}

type budgetOverviewResponse struct {
	ID          int                        `json:"id"`
	BudgetMonth string                     `json:"budget_month"`
	TotalAmount float64                    `json:"total_amount"`
	TotalSpent  float64                    `json:"total_spent"`
	Remaining   float64                    `json:"remaining"`
	Categories  []categorySpendingResponse `json:"categories"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   *time.Time                 `json:"updated_at"`
}

type budgetListData struct {
	Budgets    []monthlyBudgetResponse `json:"budgets"`
	Total      int                     `json:"total"`
	Page       int                     `json:"page"`
	PageSize   int                     `json:"page_size"`
	TotalPages int                     `json:"total_pages"`
}

func toMonthlyBudgetResponse(b *models.MonthlyBudget) monthlyBudgetResponse {
	return monthlyBudgetResponse{
		ID:          b.ID,
		BudgetMonth: b.BudgetMonth.Format("2006-01"),
		TotalAmount: b.TotalAmount,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func toBudgetOverviewResponse(result *services.BudgetOverviewResult) budgetOverviewResponse {
	cats := make([]categorySpendingResponse, 0, len(result.Categories))
	for _, c := range result.Categories {
		cats = append(cats, categorySpendingResponse{
			ID:              c.ID,
			CategoryID:      c.CategoryID,
			CategoryName:    c.CategoryName,
			AllocatedAmount: c.AllocatedAmount,
			SpentAmount:     c.SpentAmount,
			Remaining:       c.AllocatedAmount - c.SpentAmount,
		})
	}
	return budgetOverviewResponse{
		ID:          result.Budget.ID,
		BudgetMonth: result.Budget.BudgetMonth.Format("2006-01"),
		TotalAmount: result.Budget.TotalAmount,
		TotalSpent:  result.TotalSpent,
		Remaining:   result.Remaining,
		Categories:  cats,
		CreatedAt:   result.Budget.CreatedAt,
		UpdatedAt:   result.Budget.UpdatedAt,
	}
}

// ── handlers ─────────────────────────────────────────────────────────────────

func (h *Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	budgetMonth, err := parseMonth(req.BudgetMonth)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "budget_month must be in YYYY-MM format")
		return
	}

	budget, err := h.services.Budgets.CreateBudget(claims.UserID, budgetMonth, req.TotalAmount)
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type resp struct {
		Success bool                  `json:"success"`
		Message string                `json:"message"`
		Data    monthlyBudgetResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusCreated, resp{
		Success: true,
		Message: "Budget created successfully",
		Data:    toMonthlyBudgetResponse(budget),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetMonth, err := parseMonth(chi.URLParam(r, "month"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	result, err := h.services.Budgets.GetBudget(claims.UserID, budgetMonth)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	type resp struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    budgetOverviewResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "ok",
		Data:    toBudgetOverviewResponse(result),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetMonth, err := parseMonth(chi.URLParam(r, "month"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	var req updateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	budget, err := h.services.Budgets.UpdateBudgetTotal(claims.UserID, budgetMonth, req.TotalAmount)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type resp struct {
		Success bool                  `json:"success"`
		Message string                `json:"message"`
		Data    monthlyBudgetResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "Budget updated successfully",
		Data:    toMonthlyBudgetResponse(budget),
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetMonth, err := parseMonth(chi.URLParam(r, "month"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	err = h.services.Budgets.DeleteBudget(claims.UserID, budgetMonth)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJsonResponse(w, http.StatusOK, ApiResponse{
		Success: true,
		Message: "Budget deleted successfully",
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) ListBudgets(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	q := r.URL.Query()
	f := models.BudgetFilter{
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
	if yr := q.Get("year"); yr != "" {
		if v, err := strconv.Atoi(yr); err == nil && v > 0 {
			f.Year = &v
		}
	}

	budgets, total, err := h.services.Budgets.ListBudgets(claims.UserID, f)
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	items := make([]monthlyBudgetResponse, 0, len(budgets))
	for i := range budgets {
		items = append(items, toMonthlyBudgetResponse(&budgets[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(f.PageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	type resp struct {
		Success bool           `json:"success"`
		Message string         `json:"message"`
		Data    budgetListData `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "ok",
		Data: budgetListData{
			Budgets:    items,
			Total:      total,
			Page:       f.Page,
			PageSize:   f.PageSize,
			TotalPages: totalPages,
		},
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) UpsertCategoryBudget(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetMonth, err := parseMonth(chi.URLParam(r, "month"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	var req upsertCategoryBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	cat, err := h.services.Budgets.UpsertCategoryBudget(claims.UserID, budgetMonth, req.CategoryID, req.AllocatedAmount)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type resp struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message"`
		Data    categoryBudgetResponse `json:"data"`
	}
	if err := writeJsonResponse(w, http.StatusOK, resp{
		Success: true,
		Message: "Category budget set successfully",
		Data: categoryBudgetResponse{
			ID:              cat.ID,
			CategoryID:      cat.CategoryID,
			AllocatedAmount: cat.AllocatedAmount,
		},
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) DeleteCategoryBudget(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetMonth, err := parseMonth(chi.URLParam(r, "month"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	categoryID, err := strconv.Atoi(chi.URLParam(r, "category_id"))
	if err != nil || categoryID <= 0 {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid category_id")
		return
	}

	err = h.services.Budgets.DeleteCategoryBudget(claims.UserID, budgetMonth, categoryID)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJsonResponse(w, http.StatusOK, ApiResponse{
		Success: true,
		Message: "Category budget removed successfully",
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *Handler) BulkSetCategoryBudgets(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromRequest(r)
	if claims == nil {
		_ = writeJsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	budgetMonth, err := parseMonth(chi.URLParam(r, "month"))
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "month must be in YYYY-MM format")
		return
	}

	var req bulkSetCategoriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	items := make([]models.BulkCategoryItem, 0, len(req.Categories))
	for _, c := range req.Categories {
		items = append(items, models.BulkCategoryItem{
			CategoryID:      c.CategoryID,
			AllocatedAmount: c.AllocatedAmount,
		})
	}

	err = h.services.Budgets.BulkSetCategoryBudgets(claims.UserID, budgetMonth, items)
	if errors.Is(err, services.ErrNotFound) {
		_ = writeJsonError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		_ = writeJsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := writeJsonResponse(w, http.StatusOK, ApiResponse{
		Success: true,
		Message: "Category budgets updated successfully",
	}); err != nil {
		_ = writeJsonError(w, http.StatusInternalServerError, "internal server error")
	}
}
