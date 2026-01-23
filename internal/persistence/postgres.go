package persistence

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresSnapshotRepository struct {
	conn *pgx.Conn
	name string
}

func NewPostgresSnapshotRepository(conn *pgx.Conn, name string) *PostgresSnapshotRepository {
	return &PostgresSnapshotRepository{
		conn: conn,
		name: name,
	}
}

func (r *PostgresSnapshotRepository) Save(
	ctx context.Context,
	data map[string][]byte,
) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, fmt.Sprintf(`DELETE FROM %s`, r.name))
	if err != nil {
		return err
	}

	for key, value := range data {
		_, err := tx.Exec(
			ctx,
			fmt.Sprintf(`INSERT INTO %s (key, value) VALUES ($1, $2)`, r.name),
			key,
			value,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresSnapshotRepository) Load(
	ctx context.Context,
) (map[string][]byte, error) {

	rows, err := r.conn.Query(
		ctx,
		fmt.Sprintf(`SELECT key, value FROM %s`, r.name),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]byte)

	for rows.Next() {
		var key string
		var value []byte

		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}

		v := make([]byte, len(value))
		copy(v, value)

		result[key] = v
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func CreateSnapshotTable(ctx context.Context, conn *pgx.Conn, name string) error {
	sqlQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		key TEXT PRIMARY KEY,
		value BYTEA NOT NULL
	);
	`, name)
	_, err := conn.Exec(ctx, sqlQuery)
	return err
}

func DropSnapshotTable(ctx context.Context, conn *pgx.Conn, name string) error {
	sqlQuery := fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, name)
	_, err := conn.Exec(ctx, sqlQuery)
	return err
}
