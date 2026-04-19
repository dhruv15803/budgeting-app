package handlers

import (
	"net/http"
)

func (h *Handler) ListExpenseCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.services.ExpenseCategories.ListAll()
	if err != nil {
		writeJsonError(w, http.StatusInternalServerError, "failed to fetch categories")
		return
	}

	type response struct {
		Success    bool        `json:"success"`
		Message    string      `json:"message"`
		Categories interface{} `json:"categories"`
	}

	writeJsonResponse(w, http.StatusOK, response{
		Success:    true,
		Message:    "categories fetched successfully",
		Categories: categories,
	})
}
