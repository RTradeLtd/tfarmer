package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/RTradeLtd/gorm"

	"github.com/RTradeLtd/cmd"
	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database"
	"github.com/RTradeLtd/database/models"
)

// Version denotes the tag of this build
var Version string

// globals
var (
	ctx    context.Context
	cancel context.CancelFunc
)

// command-line flags
var (
	devMode    *bool
	debug      *bool
	configPath *string
	dbNoSSL    *bool
	dbMigrate  *bool
	grpcNoSSL  *bool
	apiPort    *string

	// bucket flags
	bucketLocation *string
)

func baseFlagSet() *flag.FlagSet {
	var f = flag.NewFlagSet("", flag.ExitOnError)

	// basic flags
	devMode = f.Bool("dev", false,
		"toggle dev mode")
	debug = f.Bool("debug", false,
		"toggle debug mode")
	configPath = f.String("config", os.Getenv("CONFIG_DAG"),
		"path to Temporal configuration")

	// db configuration
	dbNoSSL = f.Bool("db.no_ssl", false,
		"toggle SSL connection with database")
	dbMigrate = f.Bool("db.migrate", false,
		"toggle whether a database migration should occur")

	// grpc configuration
	grpcNoSSL = f.Bool("grpc.no_ssl", false,
		"toggle SSL connection with GRPC services")

	// api configuration
	apiPort = f.String("api.port", "6767",
		"set port to expose API on")

	return f
}

func newDB(cfg config.TemporalConfig, noSSL bool) (*gorm.DB, error) {
	return database.OpenDBConnection(database.DBOptions{
		User:           cfg.Database.Username,
		Password:       cfg.Database.Password,
		Address:        cfg.Database.URL,
		Port:           cfg.Database.Port,
		SSLModeDisable: noSSL,
	})
}

var commands = map[string]cmd.Cmd{
	"user": {
		Hidden:      true,
		Blurb:       "create a user",
		Description: "Create a Temporal user. Provide args as username, password, email. Do not use in production.",
		Args:        []string{"user", "pass", "email"},
		Action: func(cfg config.TemporalConfig, args map[string]string) {
			fmt.Printf("creating user '%s' (%s)...\n", args["user"], args["email"])
			d, err := database.Initialize(&cfg, database.Options{
				SSLModeDisable: *dbNoSSL,
				RunMigrations:  *dbMigrate,
			})
			if err != nil {
				fmt.Println("failed to initialize database connection", err)
				os.Exit(1)
			}
			// create user account
			if _, err := models.NewUserManager(d.DB).NewUserAccount(
				args["user"], args["pass"], args["email"],
			); err != nil {
				fmt.Println("failed to create user account", err)
				os.Exit(1)
			}
			// add credits
			if _, err := models.NewUserManager(d.DB).AddCredits(args["user"], 99999999); err != nil {
				fmt.Println("failed to grant credits to user account", err)
				os.Exit(1)
			}
			// generate email activation token
			userModel, err := models.NewUserManager(d.DB).GenerateEmailVerificationToken(args["user"])
			if err != nil {
				fmt.Println("failed to generate email verification token", err)
				os.Exit(1)
			}
			// activate email
			if _, err := models.NewUserManager(d.DB).ValidateEmailVerificationToken(args["user"], userModel.EmailVerificationToken); err != nil {
				fmt.Println("failed to activate email", err)
				os.Exit(1)
			}
		},
	},
}

func main() {
	if Version == "" {
		Version = "latest"
	}

	// initialize global context
	ctx, cancel = context.WithCancel(context.Background())

	// create app
	tfarmer := cmd.New(commands, cmd.Config{
		Name:     "Temporal Farmer",
		ExecName: "tfarmer",
		Version:  Version,
		Desc:     "used to scrape data from temporal's databases",
		Options:  baseFlagSet(),
	})

	// run no-config commands, exit if command was run
	if exit := tfarmer.PreRun(nil, os.Args[1:]); exit == cmd.CodeOK {
		os.Exit(0)
	}

	// load config
	tCfg, err := config.LoadConfig(*configPath)
	if err != nil {
		println("failed to load config at", *configPath)
		os.Exit(1)
	}

	// load arguments
	flags := map[string]string{
		"dbPass":  tCfg.Database.Password,
		"dbURL":   tCfg.Database.URL,
		"dbUser":  tCfg.Database.Username,
		"version": Version,
	}

	// execute
	os.Exit(tfarmer.Run(*tCfg, flags, os.Args[1:]))
}
