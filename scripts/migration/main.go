package main

import (
	"bytes"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Config struct {
	Driver    *string
	Host      *string
	Port      *int
	Username  *string
	Password  *string
	DBName    *string
	DDLFolder *string
	SeedsFile *string
}

func getConf() *Config {
	driver := flag.String("driver", "postgres", "DEFAULT TO 'postgres' database driver (only postgres is supported)")
	db := os.Getenv("SSQL_DB")
	host := os.Getenv("SSQL_HOST")
	portStr := os.Getenv("SSQL_PORT")
	port, _ := strconv.Atoi(portStr)
	username := os.Getenv("SSQL_USERNAME")
	password := os.Getenv("SSQL_PASSWORD")
	ddl := os.Getenv("SSQL_DDL")
	seeds := os.Getenv("SSQL_SEEDS")
	flag.Parse()
	config := Config{
		Driver:    driver,
		Host:      &host,
		Port:      &port,
		Username:  &username,
		Password:  &password,
		DBName:    &db,
		DDLFolder: &ddl,
		SeedsFile: &seeds,
	}
	return &config
}

func main() {
	config := getConf()
	if *config.Driver != "postgres" {
		panic("database driver not supported, only 'postgres' supported at the moment")
	}
	log.Info().Interface("config", config).Msg("config")
	ctx := context.Background()
	// Listening for commands.
	action := ""
	steps := 1
	if len(os.Args) > 1 {
		action = os.Args[1]
		fmt.Print("Command: " + action + " ")
	}
	fmt.Println("")
	// Another DB connection for migrations only.
	if action == "createdb" || action == "deletedb" {
		fmt.Printf("Executing '%s'...\n", action)
		runCmd(ctx, config, action)
		return
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		*config.Host,
		*config.Port,
		*config.Username,
		*config.Password,
		*config.DBName)
	log.Info().Msgf("Postgres Config %s", psqlInfo)
	dbM, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error().Err(err).Msg("sql.Open")
		log.Info().Msgf("cmd: %s", psqlInfo)
		panic(err)
	}
	defer dbM.Close()
	driver, err := postgres.WithInstance(dbM, &postgres.Config{})
	if err != nil {
		log.Error().Err(err).Msg("postgres.WithInstance")
		log.Info().Msgf("cmd: %s", psqlInfo)
		panic(err)
	}
	dbPath := "file://" + *config.DDLFolder
	m, err := migrate.NewWithDatabaseInstance(
		dbPath,
		"postgres", driver)
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("DDL folder: %s", dbPath)
	defer m.Close()
	currentVersion, _, _ := m.Version()
	fmt.Println("Current Version: " + fmt.Sprint(currentVersion))

	if action == "bootstrap" {
		fmt.Println("Flushing everything first and then running all migrations and seeds...")
		if qErr := m.Down(); qErr != nil {
			log.Error().Err(qErr).Msgf("Error downgrading DB")
		}
		if qErr := m.Up(); qErr != nil {
			log.Error().Err(qErr).Msgf("Error upgrading DB")
		}
		runCmd(ctx, config, "seeds")
		fmt.Println("Done!")
	} else if action == "up" {
		fmt.Println("Migrating up " + strconv.Itoa(steps) + " step(s)...")
		m.Steps(steps)
		fmt.Println("Done!")
	} else if action == "down" {
		fmt.Println("Migrating down " + strconv.Itoa(steps) + " step(s)...")
		m.Steps(-steps)
		fmt.Println("Done!")
	} else if action == "flush" || action == "reset" {
		fmt.Println("Flushing database...")
		if qErr := m.Down(); qErr != nil {
			log.Error().Err(qErr).Msgf("Error flushing DB")
		}
		fmt.Println("Done!")
	} else if action == "seed" || action == "seeds" {
		fmt.Println("Running seeds...")
		runCmd(ctx, config, "seeds")
	}
	newVersion, _, _ := m.Version()
	fmt.Println("New Version: " + fmt.Sprint(newVersion))
}

func runCmd(ctx context.Context, config *Config, op string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var cmd *exec.Cmd
	if op == "seeds" {
		cmd = exec.Command("psql", "-U", *config.Username, "-h", *config.Host, "-d", *config.DBName, "-a", "-f", *config.SeedsFile)
	} else if op == "createdb" {
		cmd = exec.Command("psql", "-U", *config.Username, "-h", *config.Host, "-a", "-f", dir+"/scripts/migration/db-create.sql")
	} else if op == "deletedb" {
		cmd = exec.Command("psql", "-U", *config.Username, "-h", *config.Host, "-a", "-f", dir+"/scripts/migration/db-delete.sql")
	} else {
		return
	}
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Panic().Msgf("Error executing query. Command Output: %+v\n: %+v, %v", out.String(), stderr.String(), err)
	}
}