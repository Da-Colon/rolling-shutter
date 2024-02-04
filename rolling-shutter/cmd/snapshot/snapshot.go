package snapshot

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/shutter-network/rolling-shutter/rolling-shutter/medley/configuration/command"
	"github.com/shutter-network/rolling-shutter/rolling-shutter/medley/db"
	"github.com/shutter-network/rolling-shutter/rolling-shutter/medley/service"
	"github.com/shutter-network/rolling-shutter/rolling-shutter/snapshot"
	"github.com/shutter-network/rolling-shutter/rolling-shutter/snapshot/database"
)

var (
	cfgFile    string
	outputFile string
)

func Cmd() *cobra.Command {
	builder := command.Build(
		main,
		command.Usage(
			"Run the Snapshot Hub communication module",
			// TODO long usage description
			"",
		),
		command.WithGenerateConfigSubcommand(),
		command.WithDumpConfigSubcommand(),
	)
	builder.AddInitDBCommand(initDB)
	return builder.Command()
}

func initDB(cfg *snapshot.Config) error {
	ctx := context.Background()
	dbpool, err := pgxpool.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "failed to connect to database")
	}
	defer dbpool.Close()
	// Ensure meta_inf table exists before any validation or creation
	createMetaInfSQL := `
        CREATE TABLE IF NOT EXISTS meta_inf (
            key VARCHAR PRIMARY KEY,
            value VARCHAR NOT NULL
        );
    `
	dbpool.Exec(ctx, createMetaInfSQL)
	return db.InitDB(ctx, dbpool, database.Definition.Name(), database.Definition)
}

func main(config *snapshot.Config) error {
	snp, err := snapshot.New(config)
	if err != nil {
		return err
	}
	return service.RunWithSighandler(context.Background(), snp)
}
