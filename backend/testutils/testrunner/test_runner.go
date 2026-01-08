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

// Run all the tests in the test file.
func (tr TestRunner) Run(m *testing.M) {
	defer tr.cleanUp()
	exitCode := m.Run()
	os.Exit(exitCode)
}

// Create and run a test.
func (tr TestRunner) Test(f func()) {
	defer tr.ClearDbData()
	f()
}

func (tr TestRunner) TestServer() *testserver.TestServer {
	return tr.ts
}

func (tr TestRunner) Context() context.Context {
	return tr.ctx
}

func (tr TestRunner) DB() *pgx.Conn {
	return tr.db
}

func (ts TestRunner) cleanUp() {
	ts.ClearDbData()
	ts.db.Close(ts.ctx)
}

func (ts TestRunner) ClearDbData() error {
	_, err := ts.db.Exec(
		ts.ctx,
		`
		DO $$ 
		DECLARE 
		    r RECORD;
		BEGIN
		    -- Only loop if there are actually tables in the public schema
		    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations') LOOP
			EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' RESTART IDENTITY CASCADE';
		    END LOOP;
		END $$;
		`,
	)

	return err
}
