package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go-slack/database"
	"go-slack/migrations"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

func main() {
	db := database.Connect()
	defer db.Close()

	ctx := context.Background()

	err := run(db, ctx, os.Args[1:])

	if err != nil {
		fmt.Println(err)
	}
}

func run(db *bun.DB, ctx context.Context, args []string) error {
	migrator := migrate.NewMigrator(db, migrations.Migrations)

	if len(args) == 0 {
		return fmt.Errorf("Migration command required (init, up, down, status, new)")
	}

	cmd := args[0]
	switch cmd {
	case "init":
		return initDb(migrator)
	case "up":
		return up(migrator)
	case "down":
		return down(migrator)
	case "status":
		return status(migrator)
	case "new":
		return createMigration(migrator, ctx, args[1:])
	default:
		return fmt.Errorf("Unknown \"%s\" migration command", cmd)
	}
}

func initDb(migrator *migrate.Migrator) error {
	fmt.Print("Creating bun migration tables...")
	err := migrator.Init(context.Background())

	if err == nil {
		fmt.Println("done")
	}

	return err
}

func up(migrator *migrate.Migrator) error {
	group, err := migrator.Migrate(context.Background())

	if err != nil {
		return err
	}

	fmt.Printf("Migrated to %s\n", group)

	return nil
}

func down(migrator *migrate.Migrator) error {
	group, err := migrator.Rollback(context.Background())

	if err != nil {
		return err
	}

	fmt.Printf("Rolled back %s\n", group)

	return nil
}

func status(migrator *migrate.Migrator) error {
	migrations, err := migrator.MigrationsWithStatus(context.Background())

	if err != nil {
		return err
	}

	unappliedMigrations := make([]migrate.Migration, 0, len(migrations))
	for _, migration := range migrations {
		if migration.IsApplied() {
			continue
		}

		unappliedMigrations = append(unappliedMigrations, migration)
	}

	if len(unappliedMigrations) == 0 {
		fmt.Println("All migrations have been applied.")
	} else {
		fmt.Println("Migrations that have not been applied:")
		for _, migration := range unappliedMigrations {
			fmt.Printf("\t%s\n", migration)
		}
	}

	return nil
}

func createMigration(migrator *migrate.Migrator, ctx context.Context, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("A name is required for the new migration.")
	}

	name := strings.Join(args, "_")

	mf, err := migrator.CreateGoMigration(ctx, name)

	if err != nil {
		return err
	}

	fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

	return nil
}
