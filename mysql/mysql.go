package mysql

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlOop struct {
	DB *sqlx.DB
}

func Init(user string, pass string, host string, dbname string, maxIdleConns, maxOpenConns, connMaxLifetime, connMaxIdleTime int) (*MysqlOop, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	connectionString := fmt.Sprintf(`%s:%s@(%s)/%s?parseTime=true`, user, pass, host, dbname)
	db, err := sqlx.ConnectContext(ctx, "mysql", connectionString)

	if err == nil {
		if maxIdleConns > 0 {
			db.SetMaxIdleConns(maxIdleConns)
		}

		if maxOpenConns > 0 {
			db.SetMaxOpenConns(maxOpenConns)
		}
	}

	// valid value: 1:minute, 2:hour 25:second
	if strconv.Itoa(connMaxLifetime) != "" {
		// set based on this
		split := strings.Split(strconv.Itoa(connMaxLifetime), ":")
		if len(split) == 2 {
			value := split[0]
			timeType := split[1]

			valueInt, err := strconv.Atoi(value)
			if err == nil && valueInt > 0 {
				// if value is int (otherwise not valid)
				if timeType == "hour" {
					db.SetConnMaxLifetime(time.Duration(valueInt) * time.Hour)
				} else if timeType == "minute" {
					db.SetConnMaxLifetime(time.Duration(valueInt) * time.Minute)
				} else if timeType == "second" {
					db.SetConnMaxLifetime(time.Duration(valueInt) * time.Second)
				}
			}
		}
	}

	mysqlClient := &MysqlOop{
		DB: db,
	}

	return mysqlClient, err
}

func (r *MysqlOop) Select(queryStatement string) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := r.DB.QueryxContext(ctx, queryStatement)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		row := make(map[string]interface{})
		if err := rows.MapScan(row); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		results = append(results, row)
	}

	return results, nil
}

func (r *MysqlOop) Update(queryStatement string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, queryStatement)
	if err != nil {
		return 0, fmt.Errorf("update query failed: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), nil
}

func (r *MysqlOop) Delete(queryStatement string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, queryStatement)
	if err != nil {
		return 0, fmt.Errorf("delete query failed: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to del row affected: %w", err)
	}

	return int(rowsAffected), nil
}

func (r *MysqlOop) Insert(queryStatement string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.DB.ExecContext(ctx, queryStatement)
	if err != nil {
		return 0, fmt.Errorf("insert query failed: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to insert row affected: %w", err)
	}

	return int(rowsAffected), nil
}
