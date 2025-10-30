package repository

import (
	"context"
	"database/sql"
	"errors"
	"jobProject/internal/model"
	"log"
	"time"
)

type SubsRepository interface {
	CreateSubRepo(ctx context.Context, model model.SubscriptionDB) error
	ReadSubRepo(ctx context.Context, id int) (model.SubscriptionDB, error)
	PatchSubByID(ctx context.Context, id int, s model.Subscription) error
}

type PostgresSubs struct {
	DB *sql.DB
}

func (r *PostgresSubs) CreateSubRepo(ctx context.Context, s model.SubscriptionDB) error {
	rows, err := r.DB.ExecContext(ctx, `INSERT INTO subs_table (service, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5)`, s.Service, s.Price, s.UserID, s.StartDate, s.EndDate)
	if err != nil {
		log.Printf("insert error: %v", err)
		return err
	}
	n, _ := rows.RowsAffected()
	log.Printf("inserted rows: %d", n)
	return nil
}

func (r *PostgresSubs) ReadSubRepo(ctx context.Context, id int) (model.SubscriptionDB, error) {
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

func (r *PostgresSubs) PatchSubByID(ctx context.Context, id int, s model.Subscription) error {
	const q = `SELECT id, service, price, user_id, start_date, end_date FROM subs_table WHERE id = $1`
	var old model.SubscriptionDB
	err := r.DB.QueryRowContext(ctx, q, id).Scan(
		&old.ID, &old.Service, &old.Price, &old.UserID, &old.StartDate, &old.EndDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return sql.ErrNoRows
	}
	var timeS *time.Time
	var timeE *time.Time
	if s.Service == nil {
		s.Service = &old.Service
	}
	if s.Price == nil {
		s.Price = &old.Price
	}
	if s.UserID == nil {
		s.UserID = &old.UserID
	}
	if s.StartDate == nil {
		timeS = &old.StartDate
	}
	if s.EndDate == nil {
		timeE = old.EndDate
	}
	const q1 = `UPDATE subs_table SET service = $1, price = $2, user_id = $3, start_date = $4, end_date = $5 WHERE id = $6`
	_, err = r.DB.ExecContext(ctx, q1, *s.Service, *s.Price, *s.UserID, timeS, timeE, id)
	if err != nil {
		return err
	}

	return nil
}
