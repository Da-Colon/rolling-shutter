// Package kprdb contains the sqlc generated files for interacting with the keyper's database
// schema.
package kprdb

import (
	"context"
	_ "embed"
	"log"
	"regexp"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

//go:embed schema.sql
// CreateKeyperTables contains the SQL statements to create the keyper namespace and tables.
var CreateKeyperTables string

func expectedSchemaVersion() string {
	rx := "-- schema-version: ([0-9]+) --"
	matches := regexp.MustCompile(rx).FindStringSubmatch(CreateKeyperTables)
	if len(matches) != 2 {
		log.Fatalf("internal error: kprdb/schema.sql is wrongly formatted, cannot find regular expression %s", rx)
	}
	return matches[1]
}

// schemaVersion is used to check that we use the right schema.
var schemaVersion = expectedSchemaVersion()

// InitKeyperDB initializes the database of the keyper. It is assumed that the db is empty.
func InitKeyperDB(ctx context.Context, dbpool *pgxpool.Pool) error {
	_, err := dbpool.Exec(ctx, CreateKeyperTables)
	if err != nil {
		return errors.Wrap(err, "failed to create keyper tables")
	}
	err = New(dbpool).InsertMeta(ctx, InsertMetaParams{Key: "version", Value: "1"})
	if err != nil {
		return errors.Wrap(err, "failed to set version in meta_inf table")
	}
	return nil
}

// ValidateKeyperDB checks that all expected tables exist in the database. If not, it returns an
// error.
func ValidateKeyperDB(ctx context.Context, dbpool *pgxpool.Pool) error {
	m, err := New(dbpool).GetMeta(ctx, "version")
	if err != nil {
		return errors.Wrap(err, "failed to get version from meta_inf table")
	}
	if m.Value != schemaVersion {
		return errors.Errorf("database has wrong schema version: expected %s, got %s", schemaVersion, m.Value)
	}
	return nil
}