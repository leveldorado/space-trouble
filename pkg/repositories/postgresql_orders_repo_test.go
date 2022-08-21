package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/leveldorado/space-trouble/pkg/tools/logger"
	"github.com/leveldorado/space-trouble/pkg/types"
	"github.com/stretchr/testify/require"
)

func prepareOrdersRepo(t *testing.T) *PostgreSQLOrdersRepo {
	url := os.Getenv("POSTGRESQL_URL")
	conn, err := GetPostgresqlConn(url)
	require.NoError(t, err)
	repo := NewPostgreSQLOrdersRepo(conn, logger.New())
	require.NoError(t, repo.CreateTables(context.TODO()))
	return repo
}

func TestPostgreSQLOrdersRepo_Insert(t *testing.T) {
	repo := prepareOrdersRepo(t)

	doc := types.Order{
		ID:            uuid.New().String(),
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Gender:        gofakeit.Gender(),
		BirthdayYear:  2000,
		BirthdayMonth: 10,
		BirthdayDay:   2,
		LaunchpadID:   uuid.New().String(),
		DestinationID: uuid.New().String(),
		LaunchDate:    gofakeit.Date(),
	}
	require.NoError(t, repo.Insert(context.TODO(), doc))

	doc.ID = uuid.New().String()
	doc.LaunchDate = gofakeit.Date()
	doc.LaunchpadID = uuid.New().String()

	require.NoError(t, repo.Insert(context.TODO(), doc))

}

func TestPostgreSQLOrdersRepo_Get(t *testing.T) {
	repo := prepareOrdersRepo(t)
	doc := types.Order{
		ID:            uuid.New().String(),
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Gender:        gofakeit.Gender(),
		BirthdayYear:  2000,
		BirthdayMonth: 10,
		BirthdayDay:   2,
		LaunchpadID:   uuid.New().String(),
		DestinationID: uuid.New().String(),
		LaunchDate:    gofakeit.Date(),
	}
	require.NoError(t, repo.Insert(context.TODO(), doc))
	fromDB, err := repo.Get(context.TODO(), doc.ID)
	require.NoError(t, err)
	fromDB.LaunchDate = fromDB.LaunchDate.UTC().Truncate(time.Millisecond)
	doc.LaunchDate = doc.LaunchDate.UTC().Truncate(time.Millisecond)
	fromDB.CreatedAt = fromDB.CreatedAt.UTC().Truncate(time.Millisecond)
	doc.CreatedAt = doc.CreatedAt.UTC().Truncate(time.Millisecond)
	require.Equal(t, doc, fromDB)
}

func TestPostgreSQLOrdersRepo_List(t *testing.T) {
	repo := prepareOrdersRepo(t)
	q := `truncate table "` + orderTableName + `";`
	_, err := repo.conn.Exec(q)
	require.NoError(t, err)
	doc := types.Order{
		ID:            uuid.New().String(),
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		Gender:        gofakeit.Gender(),
		BirthdayYear:  2000,
		BirthdayMonth: 10,
		BirthdayDay:   2,
		LaunchpadID:   uuid.New().String(),
		DestinationID: uuid.New().String(),
		LaunchDate:    gofakeit.Date(),
		CreatedAt:     time.Now().UTC(),
	}
	require.NoError(t, repo.Insert(context.TODO(), doc))
	list, err := repo.List(context.TODO(), 10, 0)
	require.NoError(t, err)
	require.Len(t, list, 1)
	list[0].LaunchDate = list[0].LaunchDate.UTC().Truncate(time.Millisecond)
	doc.LaunchDate = doc.LaunchDate.UTC().Truncate(time.Millisecond)
	list[0].CreatedAt = list[0].CreatedAt.UTC().Truncate(time.Millisecond)
	doc.CreatedAt = doc.CreatedAt.UTC().Truncate(time.Millisecond)

	require.Equal(t, []types.Order{doc}, list)
}

func TestPostgreSQLOrdersRepo_Delete(t *testing.T) {
	repo := prepareOrdersRepo(t)
	doc := types.Order{
		ID:            uuid.New().String(),
		DestinationID: uuid.New().String(),
	}
	require.NoError(t, repo.Insert(context.TODO(), doc))
	require.NoError(t, repo.Delete(context.TODO(), doc.ID))
	_, err := repo.Get(context.TODO(), doc.ID)
	require.Error(t, err)
	require.True(t, errors.As(err, &types.ErrNotFound{}))
	require.NoError(t, repo.Delete(context.TODO(), doc.ID))
}
