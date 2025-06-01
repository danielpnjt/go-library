package postgresql

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOop struct {
	DB *pgxpool.Pool
}

func Init(user, pass, host, dbname string, maxConns int, connMaxLifetime, connMaxIdleTime string) (*PostgresOop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, pass, host, dbname)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	if maxConns > 0 {
		config.MaxConns = int32(maxConns)
	}

	if connMaxLifetime != "" {
		if dur := ParseDuration(connMaxLifetime); dur > 0 {
			config.MaxConnLifetime = dur
		}
	}

	if connMaxIdleTime != "" {
		if dur := ParseDuration(connMaxIdleTime); dur > 0 {
			config.MaxConnIdleTime = dur
		}
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &PostgresOop{
		DB: dbpool,
	}, nil
}

func ParseDuration(input string) time.Duration {
	parts := strings.Split(input, ":")
	if len(parts) != 2 {
		return 0
	}
	value, err := strconv.Atoi(parts[0])
	if err != nil || value <= 0 {
		return 0
	}
	switch parts[1] {
	case "second":
		return time.Duration(value) * time.Second
	case "minute":
		return time.Duration(value) * time.Minute
	case "hour":
		return time.Duration(value) * time.Hour
	default:
		return 0
	}
}

func (r *PostgresOop) Select(queryStatement string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.DB.Query(ctx, queryStatement)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name)
	}

	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return results, nil
}

func (r *PostgresOop) Update(queryStatement string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.Exec(ctx, queryStatement)
	if err != nil {
		return 0, fmt.Errorf("update query failed: %w", err)
	}

	rowsAffected := int(result.RowsAffected())
	return rowsAffected, nil
}

func (r *PostgresOop) Delete(queryStatement string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.Exec(ctx, queryStatement)
	if err != nil {
		return 0, fmt.Errorf("delete query failed: %w", err)
	}

	rowsAffected := int(result.RowsAffected())
	return rowsAffected, nil
}

func (r *PostgresOop) Insert(queryStatement string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.Exec(ctx, queryStatement)
	if err != nil {
		return 0, fmt.Errorf("insert query failed: %w", err)
	}

	rowsAffected := int(result.RowsAffected())
	return rowsAffected, nil
}
