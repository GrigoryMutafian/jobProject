package usecase

import (
	"context"
	"errors"
	"jobProject/internal/conv"
	"jobProject/internal/model"
	"jobProject/internal/repository"
	"strings"
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

func (uc *SubUsecase) CreateColumnUC(ctx context.Context, s model.Subscription) error {
	err := validateSubscription(s)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	dbSub, err := conv.ParsedDates(s)
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	return uc.Repo.CreateColumn(ctx, dbSub)
}

func validateSubscription(s model.Subscription) error {
	if s.Price != nil && *s.Price < 0 {
		return errors.New("price must be not less then 0")
	}
	if s.Service != nil && (utf8.RuneCountInString(*s.Service) == 0 || strings.TrimSpace(*s.Service) == "") {
		return errors.New("service name is empty")
	}
	if s.UserID != nil && utf8.RuneCountInString(*s.UserID) != 36 {
		return errors.New("validate userID length error, must be 36 chars")
	}
	if s.StartDate != nil {
		err := monthYearValidate(*s.StartDate, s.EndDate)
		if err != nil {
			return errors.Join(ErrValidation, err)
		}
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

func (uc *SubUsecase) ReadColumnUC(ctx context.Context, id int) (model.SubscriptionDB, error) {
	if id <= 0 {
		return model.SubscriptionDB{}, errors.Join(ErrValidation, errors.New("id in query must be not less then 0"))
	}
	sub, err := uc.Repo.ReadColumn(ctx, id)
	if err != nil {
		return model.SubscriptionDB{}, err
	}
	return sub, nil
}

func (uc *SubUsecase) PatchColumnByID(ctx context.Context, id int, s model.Subscription) error {
	err := validateSubscription(s)
	if s.Service == nil && s.Price == nil && s.UserID == nil && s.StartDate == nil && s.EndDate == nil {
		return errors.Join(ErrValidation, errors.New("no data to update"))
	}
	if id <= 0 {
		return errors.Join(ErrValidation, errors.New("id in query must be not less then 0"))
	}
	if s.Service != nil && strings.TrimSpace(*s.Service) == "" {
		return errors.Join(ErrValidation, errors.New("service name is empty"))
	}
	if err != nil {
		return errors.Join(ErrValidation, err)
	}
	err = uc.Repo.PatchColumnByID(ctx, id, s)
	if err != nil {
		return err
	}
	return nil
}

func (uc *SubUsecase) DeleteColumnByID(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.Join(ErrValidation, errors.New("id in query must be not less then 0"))
	}
	err := uc.Repo.DeleteColumnByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *SubUsecase) TotalPriceByPeriod(ctx context.Context, userID, service string, from, to time.Time) (int, error) {
	if userID == "" || service == "" {
		return 0, errors.Join(ErrValidation, errors.New("user_id/service required"))
	}
	if from.After(to) {
		return 0, errors.Join(ErrValidation, errors.New("invalid period: from > to"))
	}

	total, err := uc.Repo.TotalPriceByPeriod(ctx, userID, service, from, to)
	if err != nil {
		return 0, err
	}
	return total, nil
}
