package testrunner

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-slack/database"
	"go-slack/testutils/testserver"
	"os"
	"testing"
)

type TestRunner struct {
	ctx context.Context
	db  *pgx.Conn
	ts  *testserver.TestServer
}

func New() (*TestRunner, error) {
	ctx := context.Background()
	db, err := connectToDb(ctx)

	if err != nil {
		return nil, err
	}

	ts, err := createTestServer(ctx, db)

	if err != nil {
		db.Close(ctx)
		return nil, err
	}

	return &TestRunner{
		ctx: ctx,
		db:  db,
		ts:  ts,
	}, nil
}

func connectToDb(ctx context.Context) (*pgx.Conn, error) {
	db, err := database.Connect(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %s", err)
	}

	return db, nil
}

func createTestServer(ctx context.Context, db *pgx.Conn) (*testserver.TestServer, error) {
	ts, err := testserver.New(ctx, db)

	if err != nil {
		db.Close(ctx)
		return nil, fmt.Errorf("Failed to create test server: %s", err)
	}

	return ts, nil
}

func (tr TestRunner) Run(m *testing.M) {
	defer tr.cleanUp()
	exitCode := m.Run()
	os.Exit(exitCode)
}

func (tr TestRunner) TestServer() *testserver.TestServer {
	return tr.ts
}

func (ts TestRunner) cleanUp() {
	ts.db.Close(ts.ctx)
}
