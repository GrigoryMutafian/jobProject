package repository

import (
	"context"
	"database/sql"
	"errors"
	"jobProject/internal/conv"
	"jobProject/internal/model"
	"log"
	"time"
)

type SubsRepository interface {
	CreateColumn(ctx context.Context, model model.SubscriptionDB) error
	ReadColumn(ctx context.Context, id int) (model.SubscriptionDB, error)
	PatchColumnByID(ctx context.Context, id int, s model.Subscription) error
	DeleteColumnByID(ctx context.Context, id int) error
	TotalPriceByPeriod(ctx context.Context, userID, service string, from, to time.Time) (int, error)
}

type PostgresSubs struct {
	DB *sql.DB
}

func (r *PostgresSubs) CreateColumn(ctx context.Context, s model.SubscriptionDB) error {
	rows, err := r.DB.ExecContext(ctx, `INSERT INTO subs_table (service, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5)`, s.Service, s.Price, s.UserID, s.StartDate, s.EndDate)
	if err != nil {
		log.Printf("insert error: %v", err)
		return err
	}
	n, _ := rows.RowsAffected()
	log.Printf("inserted rows: %d", n)
	return nil
}

func (r *PostgresSubs) ReadColumn(ctx context.Context, id int) (model.SubscriptionDB, error) {
	const q = `SELECT id, service, price, user_id, start_date, end_date FROM subs_table WHERE id = $1`
	var s model.SubscriptionDB
	err := r.DB.QueryRowContext(ctx, q, id).Scan(
		&s.ID, &s.Service, &s.Price, &s.UserID, &s.StartDate, &s.EndDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return model.SubscriptionDB{}, sql.ErrNoRows
	}
	if err != nil {
		return model.SubscriptionDB{}, err
	}
	return s, nil
}

func (r *PostgresSubs) PatchColumnByID(ctx context.Context, id int, s model.Subscription) error {
	const q = `SELECT id, service, price, user_id, start_date, end_date FROM subs_table WHERE id = $1`
	var old model.SubscriptionDB
	err := r.DB.QueryRowContext(ctx, q, id).Scan(
		&old.ID, &old.Service, &old.Price, &old.UserID, &old.StartDate, &old.EndDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows
	}

	if s.Service == nil {
		s.Service = &old.Service
	}
	if s.Price == nil {
		s.Price = &old.Price
	}
	if s.UserID == nil {
		s.UserID = &old.UserID
	}
	var timeS *time.Time
	if s.StartDate == nil && old.StartDate != (time.Time{}) {
		timeS = &old.StartDate
	} else if s.StartDate != nil {
		parsed, _ := conv.ParseMMYYYY(*s.StartDate)
		timeS = &parsed
	}
	var timeE *time.Time
	if s.EndDate == nil {
		timeE = old.EndDate
	} else if s.EndDate != nil {
		parsed, _ := conv.ParseMMYYYY(*s.EndDate)
		timeE = &parsed
	}
	const q1 = `UPDATE subs_table SET service = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`
	_, err = r.DB.ExecContext(ctx, q1, *s.Service, *s.Price, *s.UserID, timeS, timeE, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresSubs) DeleteColumnByID(ctx context.Context, id int) error {
	const q = `DELETE FROM subs_table WHERE id = $1`
	row, err := r.DB.ExecContext(ctx, q, id)
	if err != nil {
		return nil
	}
	affected, err := row.RowsAffected()
	if err != nil {
		return nil
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresSubs) TotalPriceByPeriod(ctx context.Context, userID, service string, from, to time.Time) (int, error) {
	const q = `SELECT COALESCE(SUM(price), 0) FROM subs_table WHERE user_id = $1 AND service = $2 AND start_date >= $3 AND (end_date <= $4 OR end_date IS NULL)`
	var total int
	err := r.DB.QueryRowContext(ctx, q, userID, service, from, to).Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
