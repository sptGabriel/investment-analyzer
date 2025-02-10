package pgtest

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var concurrentPool *pgxpool.Pool

func NewDB(t *testing.T, dbName string) *pgxpool.Pool {
	t.Helper()

	if dbName == "" {
		require.FailNow(t, "dbName cannot be an empty string")
	}

	const dbNameMaxLen = 32

	hash := md5.Sum([]byte(dbName))
	dbName = "test_" + hex.EncodeToString(hash[:])
	dbName = dbName[:dbNameMaxLen]

	pool := concurrentPool

	_, err := pool.Exec(context.Background(), fmt.Sprintf("drop database if exists %s;", dbName))
	require.NoError(t, err)

	_, err = pool.Exec(context.Background(), fmt.Sprintf("create database %s;", dbName))
	require.NoError(t, err)

	// replace the concurrent pool database for the new one
	connString := pool.Config().ConnString()
	index := strings.LastIndex(connString, pool.Config().ConnConfig.Database)
	connString = connString[:index] + dbName + connString[index+len(pool.Config().ConnConfig.Database):]
	connString += "&pool_min_conns=1&pool_max_conns=2"

	newPool, err := pgxpool.New(context.Background(), connString)
	require.NoError(t, err)

	t.Cleanup(func() {
		newPool.Close()

		_, err = pool.Exec(context.Background(), fmt.Sprintf("drop database %s", dbName))
		require.NoError(t, err)
	})

	return newPool
}
