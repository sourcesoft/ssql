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
		LogLevel: 4,
	}

	client, err := ssql.NewClient(ctx, dbCon, &options)
	if err != nil {
		log.Error().Err(err).Msg("Error creating postgres client")
		panic(err)
	}

	// Sample record.
	fID := "7f8d1637-ca82-4b1b-91dc-0828c98ebb34"
	fUsername := "test"
	fEmail := "test@domain.com"
	ts := 1673899847

	// Insert a new row.
	newUser := User{
		ID:        &fID,
		Username:  &fUsername,
		Email:     &fEmail,
		UpdatedAt: &ts,
		CreatedAt: &ts,
	}
	_, err = client.Insert(ctx, "user", newUser)
	if err != nil {
		log.Error().Err(err).Msg("Postgres create user error")
		panic(err)
	}
	log.Error().Err(err).Msg("insert done")

	// Update row by ID.
	fEmail = "new@test.com"
	newUser.Email = &fEmail
	res, err := client.UpdateByID(ctx, "user", "id", fID, newUser)
	if err != nil {
		log.Error().Err(err).Msg("Postgres update user error")
		panic(err)
	}
	if count, err := (*res).RowsAffected(); count < 1 {
		log.Error().Err(err).Msg("Postgres update user error, or not found")
		panic(err)
	}

	// Select by row by ID.
	// Table name: `user`
	// Column name for id: `id`
	rows, err := client.SelectByID(ctx, "user", "id", "7f8d1637-ca82-4b1b-91dc-0828c98ebb34")
	if err != nil {
		log.Error().Err(err).Msg("Cannot select by ID")
		panic(err)
	}
	var resp User
	if err := ssql.ScanOne(&resp, rows); err != nil {
		log.Error().Err(err).Msg("Cannot get resp by ID from Postgres")
		panic(err)
	}
	log.Debug().Interface("User", resp).Msg("response")

	// Delete row by ID.
	res, err = client.DeleteByID(ctx, "user", "id", fID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot delete user by ID from Postgres")
		panic(err)
	}
	if count, err := (*res).RowsAffected(); count < 1 {
		log.Error().Err(err).Msg("User not found")
		panic(err)
	}
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