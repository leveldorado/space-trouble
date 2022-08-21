package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	customerInfoTableName = "customer_info"
	orderTableName        = "order"
)

type PostgreSQLOrdersRepo struct {
	conn *sql.DB
	log  logrus.FieldLogger
}

func NewPostgreSQLOrdersRepo(conn *sql.DB, log logrus.FieldLogger) *PostgreSQLOrdersRepo {
	return &PostgreSQLOrdersRepo{conn: conn, log: log}
}

func (r *PostgreSQLOrdersRepo) CreateTables(ctx context.Context) error {
	customerTableCreateQuery := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS "%s" (
    id             uuid,
    first_name     text,
    last_name      text,
    birthday_year  int,
    birthday_month int,
    birthday_day   int,
    gender         text,
    PRIMARY KEY(id)
);
`, customerInfoTableName)
	orderTableCreateQuery := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS "%s"  (
    id             uuid,
    customer_id    uuid,
    launchpad_id   text,
    destination_id uuid,
    launch_date    timestamp with time zone,
    created_at     timestamp,
    PRIMARY KEY(id)
);
`, orderTableName)
	for _, q := range []string{customerTableCreateQuery, orderTableCreateQuery} {
		if _, err := r.conn.ExecContext(ctx, q); err != nil {
			return errors.Wrapf(err, `failed to exec query: q - %s`, q)
		}
	}
	return nil
}

func (r *PostgreSQLOrdersRepo) Insert(ctx context.Context, doc types.Order) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return errors.Wrap(err, `failed to begin transaction`)
	}
	if err = insertOrderWithTransaction(ctx, tx, doc); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			r.log.WithField("err", err.Error()).Error("failed to rollback")
		}
		return errors.Wrapf(err, `failed to insert order: doc - %+v`, doc)
	}
	return errors.Wrapf(tx.Commit(), `failed to commit: doc - %+v`, doc)
}

func insertOrderWithTransaction(ctx context.Context, tx *sql.Tx, doc types.Order) error {
	customerID, err := obtainCustomerID(ctx, tx, doc)
	if err != nil {
		return err
	}
	q := `INSERT INTO "` + orderTableName + `" (id, customer_id, launchpad_id, destination_id, launch_date, created_at) ` +
		`VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, q, doc.ID, customerID, doc.LaunchpadID, doc.DestinationID, doc.LaunchDate, doc.CreatedAt)
	return errors.Wrapf(err, `failed to exec query: q - %s, doc - %v`, q, doc)
}

func obtainCustomerID(ctx context.Context, tx *sql.Tx, doc types.Order) (string, error) {
	q := `SELECT id FROM ` + customerInfoTableName +
		` WHERE first_name = $1 AND last_name = $2 AND gender = $3 ` +
		` AND birthday_year = $4 AND birthday_month = $5 AND birthday_day = $6 `
	var id string
	err := tx.QueryRowContext(
		ctx,
		q,
		doc.FirstName,
		doc.LastName,
		doc.Gender,
		doc.BirthdayYear,
		doc.BirthdayMonth,
		doc.BirthdayDay,
	).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return "", errors.Wrapf(err, `failed to obtain id: q - %s, doc - %+v`, q, doc)
	}
	id = uuid.New().String()
	q = `INSERT INTO ` + customerInfoTableName +
		` (id, first_name, last_name, gender, birthday_year, birthday_month, birthday_day) ` +
		` VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(
		ctx,
		q,
		id,
		doc.FirstName,
		doc.LastName,
		doc.Gender,
		doc.BirthdayYear,
		doc.BirthdayMonth,
		doc.BirthdayDay,
	)
	return id, errors.Wrapf(err, `failed to insert customer info: q - %s, doc - %+v`, q, doc)
}

func (r *PostgreSQLOrdersRepo) Get(ctx context.Context, id string) (types.Order, error) {
	q := `SELECT o.id, c.first_name, c.last_name, c.gender, c.birthday_year, c.birthday_month, c.birthday_day, ` +
		`o.launchpad_id, o.destination_id, o.launch_date, o.created_at FROM "` + orderTableName + `" o JOIN ` +
		customerInfoTableName + ` c ON o.customer_id = c.id WHERE o.id = $1;`
	doc := types.Order{}
	err := r.conn.QueryRowContext(ctx, q, id).Scan(
		&doc.ID,
		&doc.FirstName,
		&doc.LastName,
		&doc.Gender,
		&doc.BirthdayYear,
		&doc.BirthdayMonth,
		&doc.BirthdayDay,
		&doc.LaunchpadID,
		&doc.DestinationID,
		&doc.LaunchDate,
		&doc.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return types.Order{}, types.ErrNotFound{}
	}
	return doc, errors.Wrapf(err, `failed to query row: id - %s, q - %s`, id, q)
}

func (r *PostgreSQLOrdersRepo) List(ctx context.Context, limit, offset int) ([]types.Order, error) {
	q := `SELECT o.id, c.first_name, c.last_name, c.gender, c.birthday_year, c.birthday_month, c.birthday_day, ` +
		`o.launchpad_id, o.destination_id, o.launch_date, o.created_at FROM "` + orderTableName + `" o JOIN ` +
		customerInfoTableName + ` c ON o.customer_id = c.id ORDER BY o.created_at ` +
		fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	rows, err := r.conn.QueryContext(ctx, q)
	if err != nil {
		return nil, errors.Wrapf(err, `failed to query rows: q - %s`, q)
	}
	var orders []types.Order
	for rows.Next() {
		doc := types.Order{}
		if err = rows.Scan(
			&doc.ID,
			&doc.FirstName,
			&doc.LastName,
			&doc.Gender,
			&doc.BirthdayYear,
			&doc.BirthdayMonth,
			&doc.BirthdayDay,
			&doc.LaunchpadID,
			&doc.DestinationID,
			&doc.LaunchDate,
			&doc.CreatedAt,
		); err != nil {
			return nil, errors.Wrapf(err, `failed to scan rows: q -  %s`, q)
		}
		orders = append(orders, doc)
	}
	return orders, errors.Wrapf(rows.Close(), `failed to close rows: q - %s`, q)
}

func (r *PostgreSQLOrdersRepo) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM "` + orderTableName + `" WHERE id = $1`
	_, err := r.conn.ExecContext(ctx, q, id)
	return errors.Wrapf(err, `failed to exec query: id - %s, q - %s`, id, q)
}
