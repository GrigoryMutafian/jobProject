package usecase

import (
	"context"
	"errors"
	"jobProject/internal/model"
	"jobProject/internal/repository"
	"time"
	"unicode/utf8"
)

var (
	ErrValidation = errors.New("validation error")
	ErrConflict   = errors.New("conflict error")
)

func IsValidationErr(err error) bool { return errors.Is(err, ErrValidation) }
func IsConflictErr(err error) bool   { return errors.Is(err, ErrConflict) }

type SubUsecase struct {
	Repo repository.SubsRepository
}

func NewSubUsecase(repo repository.SubsRepository) *SubUsecase {
	return &SubUsecase{Repo: repo}
}

func (uc *SubUsecase) Create(ctx context.Context, s model.Subscription) error {
	err := validateSubscription(s)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	err = monthYearValidate(s.StartDate, s.EndDate)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	return nil
}

func validateSubscription(s model.Subscription) error {
	if s.Price < 0 {
		return errors.New("price must be not less then 0")
	}
	if utf8.RuneCountInString(s.Service) == 0 {
		return errors.New("service name is empty")
	}
	if utf8.RuneCountInString(s.UserID) < 2 || utf8.RuneCountInString(s.UserID) > 64 {
		return errors.New("validate userID length error")
	}
	return nil
}

var ErrBadYearMonth = errors.New("invalid date format want MM-YYYY")

func monthYearValidate(start string, end *string) error {
	StartTime, err := ParseMMYYYY(start)
	if err != nil {
		return err
	}

	if end == nil {
		return nil
	}

	EndTime, err := ParseMMYYYY(*end)
	if err != nil {
		return err
	}
	if StartTime.After(EndTime) {
		return errors.New("end_date must be more then start_date")
	}
	return nil
}

func ParseMMYYYY(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, ErrBadYearMonth
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
