package conv

import (
	"jobProject/internal/model"
	"time"
)

func ParsedDates(s model.Subscription) (model.SubscriptionDB, error) {
	start, err := ParseMMYYYY(s.StartDate)
	if err != nil {
		return model.SubscriptionDB{}, err
	}

	var endConv *time.Time
	if s.EndDate != nil {
		end, err := ParseMMYYYY(*s.EndDate)
		if err != nil {
			return model.SubscriptionDB{}, err
		}
		endConv = &end
	}
	return model.SubscriptionDB{
		Service:   s.Service,
		Price:     s.Price,
		UserID:    s.UserID,
		StartDate: start,
		EndDate:   endConv,
	}, nil
}

func ParseMMYYYY(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
