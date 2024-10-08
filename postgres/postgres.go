/*
Copyright © 2024 Acronis International GmbH.

Released under MIT license.
*/

// Package postgres provides helpers for working Postgres database.
// Should be imported explicitly.
// To register postgres as retryable func use side effect import like so:
//
//	import _ "github.com/acronis/go-dbkit/postgres"
package postgres

import (
	pg "github.com/lib/pq"

	"github.com/acronis/go-dbkit"
)

// nolint
func init() {
	db.RegisterIsRetryableFunc(&pg.Driver{}, func(err error) bool {
		if pgErr, ok := err.(*pg.Error); ok {
			name := db.PostgresErrCode(pgErr.Code.Name())
			switch name {
			case db.PostgresErrCodeDeadlockDetected:
				return true
			case db.PostgresErrCodeSerializationFailure:
				return true
			}
		}
		return false
	})
}

// CheckPostgresError checks if the passed error relates to Postgres and it's internal code matches the one from the argument.
// nolint: staticcheck // lib/pq using is deprecated. Use pgx Postgres driver.
func CheckPostgresError(err error, errCode db.PostgresErrCode) bool {
	if pgErr, ok := err.(*pg.Error); ok {
		return pgErr.Code.Name() == string(errCode)
	}
	return false
}
