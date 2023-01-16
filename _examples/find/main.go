package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/sourcesoft/ssql"
)

func main() {
	ctx := context.Background()
	db := os.Getenv("SSQL_DB")
	host := os.Getenv("SSQL_HOST")
	portStr := os.Getenv("SSQL_PORT")
	port, _ := strconv.Atoi(portStr)
	username := os.Getenv("SSQL_USERNAME")
	password := os.Getenv("SSQL_PASSWORD")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host,
		port,
		username,
		password,
		db)
	log.Debug().Msgf("Postgres Config %s", psqlInfo)
	dbCon, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error().Err(err).Msg("Error creating postgres connection")
		panic(err)
	}

	options := ssql.Options{
		Tag:      "sql",
		LogLevel: ssql.LevelDebug,
	}

	client, err := ssql.NewClient(ctx, dbCon, &options)
	if err != nil {
		log.Error().Err(err).Msg("Error creating postgres client")
		panic(err)
	}

	// Preparing the query options.
	// Fields to fetch, set to null to fetch all fields (uses '*').
	dbFields := map[string]bool{
		"id":         true,
		"email":      true,
		"username":   true,
		"active":     true,
		"created_at": true,
	}
	// Add some custom conditions.
	conds := []*ssql.ConditionPair{{
		Field: "active",
		Value: true,
		Op:    "=",
	}}
	// setting up pagination.
	shouldReturnTotalCount := true
	limit := 10
	params := ssql.Params{
		OffsetParams: &ssql.OffsetParams{
			Limit: &limit,
			// There's also offset available.
		},
		// There's also 'order' available.
		// There's also 'cursor' pagination available.
	}
	opts := ssql.SQLQueryOptions{
		Table:          "user",
		Fields:         dbFields,
		WithTotalCount: shouldReturnTotalCount,
		Params:         &params,
		Conditions:     conds,
	}

	// Executing the query.
	rows, total, err := client.Find(ctx, &opts)
	if err != nil {
		log.Error().Err(err).Msg("Cannot find users")
		panic(err)
	}

	// Reading through the results.
	var users []User
	for rows.Next() {
		var user User
		if err := ssql.ScanRow(&user, rows); err != nil {
			log.Error().Err(err).Msg("Cannot scan users")
		}
		users = append(users, user)
	}
	log.Debug().Interface("users", users).Msgf("Found users successfully with total of %d", *total)
}

type User struct {
	ID            *string `json:"id,omitempty" sql:"id" graph:"id" rel:"pk"`
	Username      *string `json:"username,omitempty" sql:"username" graph:"username"` // Username handle.
	Email         *string `json:"email,omitempty" sql:"email" graph:"email"`          // Primary email.
	EmailVerified *bool   `json:"emailVerified,omitempty" sql:"email_verified" graph:"emailVerified"`
	Active        *bool   `json:"active,omitempty" sql:"active" graph:"active"`
	UpdatedAt     *int    `json:"updatedAt,omitempty" sql:"updated_at" graph:"updatedAt"`
	CreatedAt     *int    `json:"createdAt,omitempty" sql:"created_at" graph:"createdAt"`
	DeletedAt     *int    `json:"deletedAt,omitempty" sql:"deleted_at" graph:"deletedAt"`
}