package usecase

import (
	"context"
	"errors"
	"jobProject/internal/conv"
	"jobProject/internal/model"
	"jobProject/internal/repository"
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

func (uc *SubUsecase) CreateSubUC(ctx context.Context, s model.Subscription) error {
	err := validateSubscription(s)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	err = monthYearValidate(s.StartDate, s.EndDate)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	dbSub, err := conv.ParsedDates(s)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	return uc.Repo.CreateSubRepo(ctx, dbSub)
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
	StartTime, err := conv.ParseMMYYYY(start)
	if err != nil {
		return err
	}

	if end == nil {
		return nil
	}

	EndTime, err := conv.ParseMMYYYY(*end)
	if err != nil {
		return err
	}
	if StartTime.After(EndTime) {
		return errors.Join(ErrValidation, errors.New("end_date must be more then start_date"))
	}
	return nil
}

func (uc *SubUsecase) ReadSubUC(ctx context.Context, id int) (model.SubscriptionDB, error) {
	if id <= 0 {
		return model.SubscriptionDB{}, errors.Join(ErrValidation, errors.New("id in query must be not less then 0"))
	}
	sub, err := uc.Repo.ReadSubRepo(ctx, id)
	if err != nil {
		return model.SubscriptionDB{}, err
	}
	return sub, nil
}
