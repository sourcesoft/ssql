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

	// Let's run our first query.
	limit := 1 // Only getting one record so we have room for the next pagination.
	params := ssql.Params{
		// OffsetParams: &ssql.OffsetParams{
		// 	Limit: &limit,
		// },
		CursorParams: &ssql.CursorParams{
			First: &limit,
		},
	}
	opts := ssql.SQLQueryOptions{
		Table:             "user",
		Params:            &params,
		MainSortField:     "created_at",
		MainSortDirection: ssql.DirectionDesc,
	}
	result, err := client.Find(ctx, &opts)
	if err != nil {
		log.Error().Err(err).Msg("Cannot find users")
		panic(err)
	}
	var users []User
	pageInfo, err := ssql.SuperScan(&users, result)
	if err != nil {
		log.Error().Err(err).Msg("SuperScan failed")
		panic(err)
	}
	log.Debug().Interface("users", users).Msg("Found users successfully")
	log.Debug().Interface("pageInfo", pageInfo).Msg("SuperScan pageInfo struct")

	// Create cursor using created_at field
	log.Debug().Str("endCursor", *pageInfo.EndCursor).Msg("Next cursor")

	// Let's get the next rows
	paramsNext := ssql.Params{
		CursorParams: &ssql.CursorParams{
			After: pageInfo.EndCursor, // If you have the previous cursor, you can pass it here to continue the pagination.
			First: &limit,             // Get first 10 results (works like LIMIT).
		},
	}
	optsNext := ssql.SQLQueryOptions{
		Table:             "user",
		Params:            &paramsNext,
		MainSortField:     "created_at",
		MainSortDirection: ssql.DirectionDesc,
	}
	result2, err := client.Find(ctx, &optsNext)
	if err != nil {
		log.Error().Err(err).Msg("Cannot find users again")
		panic(err)
	}
	var usersNext []User
	pageInfo2, err := ssql.SuperScan(&usersNext, result2)
	if err != nil {
		log.Error().Err(err).Msg("SuperScan failed")
		panic(err)
	}
	log.Debug().Interface("usersNext", usersNext).Msg("Found users successfully")
	log.Debug().Interface("pageInfo", pageInfo2).Msg("SuperScan pageInfo struct")
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