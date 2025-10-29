package repository

import (
	"context"
	"database/sql"
	"errors"
	"jobProject/internal/model"
	"log"
)

type SubsRepository interface {
	CreateSubRepo(ctx context.Context, model model.SubscriptionDB) error
	ReadSubRepo(ctx context.Context, id int) (model.SubscriptionDB, error)
}

type PostgresSubs struct {
	DB *sql.DB
}

func (r *PostgresSubs) CreateSubRepo(ctx context.Context, model model.SubscriptionDB) error {
	rows, err := r.DB.ExecContext(ctx, `INSERT INTO subs_table (service, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5)`, model.Service, model.Price, model.UserID, model.StartDate, model.EndDate)
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
