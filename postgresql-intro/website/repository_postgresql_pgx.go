package website

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgreSQLPGXRepository struct {
	db *pgxpool.Pool
}

func NewPostgreSQLPGXRepository(db *pgxpool.Pool) *PostgreSQLPGXRepository {
	return &PostgreSQLPGXRepository{
		db: db,
	}
}

func (r *PostgreSQLPGXRepository) Migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS websites(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL,
		rank INT NOT NULL
	);
	`

	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *PostgreSQLPGXRepository) Create(ctx context.Context, website Website) (*Website, error) {
	var id int64
	err := r.db.QueryRow(ctx, "INSERT INTO websites(name, url, rank) values($1, $2, $3) RETURNING id", website.Name, website.URL, website.Rank).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}
	website.ID = id

	return &website, nil
}

func (r *PostgreSQLPGXRepository) All(ctx context.Context) ([]Website, error) {
	rows, err := r.db.Query(ctx, "SELECT * FROM websites")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Website
	for rows.Next() {
		var website Website
		if err := rows.Scan(&website.ID, &website.Name, &website.URL, &website.Rank); err != nil {
			return nil, err
		}
		all = append(all, website)
	}
	return all, nil
}

func (r *PostgreSQLPGXRepository) GetByName(ctx context.Context, name string) (*Website, error) {
	row := r.db.QueryRow(ctx, "SELECT * FROM websites WHERE name = $1", name)

	var website Website
	if err := row.Scan(&website.ID, &website.Name, &website.URL, &website.Rank); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &website, nil
}

func (r *PostgreSQLPGXRepository) Update(ctx context.Context, id int64, updated Website) (*Website, error) {
	res, err := r.db.Exec(ctx, "UPDATE websites SET name = $1, url = $2, rank = $3 WHERE id = $4", updated.Name, updated.URL, updated.Rank, id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *PostgreSQLPGXRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM websites WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
