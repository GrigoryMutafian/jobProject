package repository

import (
	"context"
	"database/sql"
	"jobProject/internal/model"
	"log"
)

type SubsRepository interface {
	Create(ctx context.Context, model model.Subscription) error
}

type PostgresSubs struct {
	DB *sql.DB
}

func (r *PostgresSubs) Create(ctx context.Context, model model.Subscription) error {
	rows, err := r.DB.ExecContext(ctx, `INSERT INTO user_subs (service_name, price, user_id, start_date, end_date) VALUES ($1,$2,$3,$4,$5)`, model.Service, model.Price, model.UserID, model.StartDate, model.EndDate)
	if err != nil {
		log.Printf("insert error: %v", err)
		return err
	}
	n, _ := rows.RowsAffected()
	log.Printf("inserted rows: %d", n)
	return nil
}
