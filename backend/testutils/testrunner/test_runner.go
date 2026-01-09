package testrunner

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go-slack/testutils/testserver"
	"os"
	"testing"
)

type TestRunner struct {
	ctx context.Context
	db  *pgxpool.Pool
	ts  *testserver.TestServer
}

func New() (*TestRunner, error) {
	ctx := context.Background()
	db, err := connectToDb(ctx)

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %w", err)
	}

	ts := testserver.New(ctx, db)

	tr := &TestRunner{
		ctx: ctx,
		db:  db,
		ts:  ts,
	}

	return tr, nil
}

func connectToDb(ctx context.Context) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(ctx, os.Getenv("DB_URL"))

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %w", err)
	}

	return db, nil
}

// Run all the tests in the test file.
func (tr TestRunner) Run(m *testing.M) {
	defer tr.cleanUp()
	exitCode := m.Run()
	os.Exit(exitCode)
}

// Create and run a test.
func (tr TestRunner) Test(test func()) {
	defer tr.ClearDbData()
	test()
}

func (tr TestRunner) TestServer() *testserver.TestServer {
	return tr.ts
}

func (tr TestRunner) Context() context.Context {
	return tr.ctx
}

func (tr TestRunner) DB() *pgxpool.Pool {
	return tr.db
}

func (ts TestRunner) cleanUp() {
	ts.ClearDbData()
	ts.db.Close()
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
