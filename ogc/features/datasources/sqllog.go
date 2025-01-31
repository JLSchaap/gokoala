package datasources

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type contextKey int

const sqlContextKey contextKey = iota

var logSQL, _ = strconv.ParseBool(os.Getenv("LOG_SQL"))

// SQLLog query logging for debugging purposes
type SQLLog struct{}

// Before callback prior to execution of the given SQL query
func (s *SQLLog) Before(ctx context.Context, _ string, _ ...any) (context.Context, error) {
	return context.WithValue(ctx, sqlContextKey, time.Now()), nil
}

// After callback once execution of the given SQL query is done
func (s *SQLLog) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	if logSQL {
		query = replaceBindVars(query, args)
		start := ctx.Value(sqlContextKey).(time.Time)

		log.Printf("\n--- SQL:\n%s\n--- SQL query took: %s\n", query, time.Since(start))
	}
	return ctx, nil
}

// replaces '?' bind vars in order to log a complete query
func replaceBindVars(query string, args []any) string {
	for _, arg := range args {
		query = strings.Replace(query, "?", fmt.Sprintf("%v", arg), 1)
	}
	return query
}
