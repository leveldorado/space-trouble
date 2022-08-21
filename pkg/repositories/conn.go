package repositories

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func GetPostgresqlConn(url string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrapf(err, `failec to open conn: url - %s`, url)
	}
	return conn, errors.Wrapf(conn.Ping(), `failed to ping: url - %s`, url)
}
