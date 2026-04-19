package services

import (
	"context"
	"fmt"
	"time"

	"github.com/dhruv15803/budgeting-app/internal/models"
	"github.com/dhruv15803/budgeting-app/internal/repositories"
)

type RecurringExpenseServiceImpl struct {
	repo   *repositories.Repository
	logger func(format string, args ...interface{})
}

func NewRecurringExpenseService(repo *repositories.Repository, logger func(format string, args ...interface{})) *RecurringExpenseServiceImpl {
	return &RecurringExpenseServiceImpl{repo: repo, logger: logger}
}

func isValidFrequency(f string) bool {
	switch f {
	case models.FrequencyDaily, models.FrequencyWeekly, models.FrequencyMonthly, models.FrequencyYearly:
		return true
	}
	return false
}

// advance returns d shifted forward by one step of the given frequency.
func advance(d time.Time, freq string) time.Time {
	switch freq {
	case models.FrequencyDaily:
		return d.AddDate(0, 0, 1)
	case models.FrequencyWeekly:
		return d.AddDate(0, 0, 7)
	case models.FrequencyMonthly:
		return d.AddDate(0, 1, 0)
	case models.FrequencyYearly:
		return d.AddDate(1, 0, 0)
	}
	return d
}

func validateRecurringInput(title string, amount float64, categoryID int, startDate time.Time, endDate *time.Time, frequency string) error {
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if categoryID <= 0 {
		return fmt.Errorf("category_id is required")
	}
	if startDate.IsZero() {
		return fmt.Errorf("start_date is required")
	}
	if !isValidFrequency(frequency) {
		return fmt.Errorf("frequency must be one of: daily, weekly, monthly, yearly")
	}
	if endDate != nil && endDate.Before(startDate) {
		return fmt.Errorf("end_date must be on or after start_date")
	}
	return nil
}

func (s *RecurringExpenseServiceImpl) CreateRecurringExpense(userID, categoryID int, title string, description *string, amount float64, startDate time.Time, endDate *time.Time, frequency string) (*models.RecurringExpense, error) {
	if err := validateRecurringInput(title, amount, categoryID, startDate, endDate, frequency); err != nil {
		return nil, err
	}
	return s.repo.RecurringExpenses.Create(userID, categoryID, title, description, amount, startDate, endDate, frequency)
}

func (s *RecurringExpenseServiceImpl) GetRecurringExpense(id, requestingUserID int) (*models.RecurringExpense, error) {
	existing, err := s.repo.RecurringExpenses.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	if existing.UserID != requestingUserID {
		return nil, ErrForbidden
	}
	return existing, nil
}

func (s *RecurringExpenseServiceImpl) UpdateRecurringExpense(id, requestingUserID, categoryID int, title string, description *string, amount float64, startDate time.Time, endDate *time.Time, frequency string, isActive bool) (*models.RecurringExpense, error) {
	if err := validateRecurringInput(title, amount, categoryID, startDate, endDate, frequency); err != nil {
		return nil, err
	}

	existing, err := s.repo.RecurringExpenses.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrNotFound
	}
	if existing.UserID != requestingUserID {
		return nil, ErrForbidden
	}

	// Keep next_occurrence managed by the generator. Only bump it forward if the user
	// pushes start_date later than the currently scheduled next_occurrence.
	nextOccurrence := existing.NextOccurrence
	if startDate.After(nextOccurrence) {
		nextOccurrence = startDate
	}

	return s.repo.RecurringExpenses.Update(id, title, description, amount, categoryID, startDate, endDate, frequency, isActive, nextOccurrence)
}

func (s *RecurringExpenseServiceImpl) DeleteRecurringExpense(id, requestingUserID int) error {
	existing, err := s.repo.RecurringExpenses.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNotFound
	}
	if existing.UserID != requestingUserID {
		return ErrForbidden
	}
	return s.repo.RecurringExpenses.Delete(id)
}

func (s *RecurringExpenseServiceImpl) ListRecurringExpenses(userID int, f models.RecurringExpenseFilter) ([]models.RecurringExpense, int, error) {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = 20
	}
	if f.PageSize > 100 {
		f.PageSize = 100
	}
	if f.Frequency != nil && *f.Frequency != "" && !isValidFrequency(*f.Frequency) {
		return nil, 0, fmt.Errorf("frequency must be one of: daily, weekly, monthly, yearly")
	}
	return s.repo.RecurringExpenses.ListByUser(userID, f)
}

// RunDueGenerator scans every active recurring expense whose next_occurrence is on
// or before today and creates the missing expense rows, advancing next_occurrence
// forward by one frequency step per generated occurrence. Each recurring row is
// processed in its own transaction so one failure never stalls the whole batch.
// Returns the total number of expenses actually inserted.
func (s *RecurringExpenseServiceImpl) RunDueGenerator(ctx context.Context) (int, error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	due, err := s.repo.RecurringExpenses.ListDue(today)
	if err != nil {
		return 0, err
	}

	generated := 0
	for i := range due {
		if err := ctx.Err(); err != nil {
			return generated, err
		}
		if n, err := s.processDueRow(&due[i], today); err != nil {
			if s.logger != nil {
				s.logger("recurring generator: row id=%d failed: %v", due[i].ID, err)
			}
		} else {
			generated += n
		}
	}
	return generated, nil
}

func (s *RecurringExpenseServiceImpl) processDueRow(r *models.RecurringExpense, today time.Time) (int, error) {
	tx, err := s.repo.BeginTx()
	if err != nil {
		return 0, err
	}
	defer func() { _ = tx.Rollback() }()

	generated := 0
	next := r.NextOccurrence

	for !next.After(today) {
		if r.EndDate != nil && next.After(*r.EndDate) {
			break
		}
		inserted, err := s.repo.Expenses.CreateFromRecurringTx(tx, r, next)
		if err != nil {
			return 0, err
		}
		if inserted {
			generated++
		}
		next = advance(next, r.Frequency)
	}

	deactivate := r.EndDate != nil && next.After(*r.EndDate)
	if err := s.repo.RecurringExpenses.AdvanceTx(tx, r.ID, next, deactivate); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return generated, nil
}
