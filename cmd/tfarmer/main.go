package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/RTradeLtd/gorm"
	"github.com/RTradeLtd/rtfs"
	"github.com/RTradeLtd/tfarmer/mail"
	"github.com/RTradeLtd/tfarmer/upload"
	"github.com/RTradeLtd/tfarmer/user"

	"github.com/RTradeLtd/cmd"
	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database"
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
	devMode        *bool
	debug          *bool
	configPath     *string
	dbNoSSL        *bool
	dbMigrate      *bool
	sendEmail      *bool
	emailRecipient *string
	recipientName  *string
	uploadType     *string
	unique         *bool
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
	unique = f.Bool("unique", false,
		"toggle whether unique checks should be performed")

	// db configuration
	dbNoSSL = f.Bool("db.no_ssl", false,
		"toggle SSL connection with database")
	dbMigrate = f.Bool("db.migrate", false,
		"toggle whether a database migration should occur")

	// email flags
	sendEmail = f.Bool("email-enabled", false,
		"used to activate email notification")
	emailRecipient = f.String("email-recipient", "",
		"email to send metrics to")
	recipientName = f.String("recipient-name", "",
		"email recipient name")

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
		Blurb:         "User based metrics",
		Description:   "Allows for gathering of user based metric information (usage, activity, etc..)",
		ChildRequired: true,
		Children: map[string]cmd.Cmd{
			"registered": {
				Blurb:       "Registered users",
				Description: "Used to get the number of registered users",
				Action: func(cfg config.TemporalConfig, args map[string]string) {
					db, err := database.Initialize(&cfg, database.Options{
						SSLModeDisable: *dbNoSSL,
						RunMigrations:  *dbMigrate,
					})
					if err != nil {
						fmt.Println("failed to initialize database connection", err.Error())
						os.Exit(1)
					}
					uf := user.NewFarmer(db.DB)
					users, err := uf.RegisteredUsers()
					if err != nil {
						fmt.Println("failed to get registered users", err.Error())
						os.Exit(1)
					}
					numberOfUsers := len(users)
					msg := fmt.Sprintf("there are %v total registered users", numberOfUsers)
					fmt.Println(msg)
					if *sendEmail {
						mm, err := mail.NewManager(&cfg, db.DB)
						if err != nil {
							fmt.Println("failed to initialize mail manager", err.Error())
							os.Exit(1)
						}
						if _, err := mm.SendEmail(
							"registered users report",
							msg,
							"text/html",
							*recipientName,
							*emailRecipient,
						); err != nil {
							fmt.Println("failed to send email report", err.Error())
							os.Exit(1)
						}
					}
				},
			},
			"free": {
				Blurb:       "Free users",
				Description: "Used to get the number of free users",
				Action: func(cfg config.TemporalConfig, args map[string]string) {
					db, err := database.Initialize(&cfg, database.Options{
						SSLModeDisable: *dbNoSSL,
						RunMigrations:  *dbMigrate,
					})
					if err != nil {
						fmt.Println("failed to initialize database connection", err.Error())
						os.Exit(1)
					}
					uf := user.NewFarmer(db.DB)
					users, err := uf.FreeUsers()
					if err != nil {
						fmt.Println("failed to get free users", err.Error())
						os.Exit(1)
					}
					numberOfUsers := len(users)
					msg := fmt.Sprintf("there are %v total free users", numberOfUsers)
					fmt.Println(msg)
					if *sendEmail {
						mm, err := mail.NewManager(&cfg, db.DB)
						if err != nil {
							fmt.Println("failed to initialize mail manager", err.Error())
							os.Exit(1)
						}
						if _, err := mm.SendEmail(
							"free users report",
							msg,
							"text/html",
							*recipientName,
							*emailRecipient,
						); err != nil {
							fmt.Println("failed to send email report", err.Error())
							os.Exit(1)
						}
					}
				},
			},
			"light": {
				Blurb:       "Light users",
				Description: "Used to get the number of light users",
				Action: func(cfg config.TemporalConfig, args map[string]string) {
					db, err := database.Initialize(&cfg, database.Options{
						SSLModeDisable: *dbNoSSL,
						RunMigrations:  *dbMigrate,
					})
					if err != nil {
						fmt.Println("failed to initialize database connection", err.Error())
						os.Exit(1)
					}
					uf := user.NewFarmer(db.DB)
					users, err := uf.LightUsers()
					if err != nil {
						fmt.Println("failed to get light users", err.Error())
						os.Exit(1)
					}
					numberOfUsers := len(users)
					msg := fmt.Sprintf("there are %v total light users", numberOfUsers)
					fmt.Println(msg)
					if *sendEmail {
						mm, err := mail.NewManager(&cfg, db.DB)
						if err != nil {
							fmt.Println("failed to initialize mail manager", err.Error())
							os.Exit(1)
						}
						if _, err := mm.SendEmail(
							"light users report",
							msg,
							"text/html",
							*recipientName,
							*emailRecipient,
						); err != nil {
							fmt.Println("failed to send email report", err.Error())
							os.Exit(1)
						}
					}
				},
			},
			"plus": {
				Blurb:       "Plus users",
				Description: "Used to get the number of plus users",
				Action: func(cfg config.TemporalConfig, args map[string]string) {
					db, err := database.Initialize(&cfg, database.Options{
						SSLModeDisable: *dbNoSSL,
						RunMigrations:  *dbMigrate,
					})
					if err != nil {
						fmt.Println("failed to initialize database connection", err.Error())
						os.Exit(1)
					}
					uf := user.NewFarmer(db.DB)
					users, err := uf.FreeUsers()
					if err != nil {
						fmt.Println("failed to get plus users", err.Error())
						os.Exit(1)
					}
					numberOfUsers := len(users)
					msg := fmt.Sprintf("there are %v total plus users", numberOfUsers)
					fmt.Println(msg)
					if *sendEmail {
						mm, err := mail.NewManager(&cfg, db.DB)
						if err != nil {
							fmt.Println("failed to initialize mail manager", err.Error())
							os.Exit(1)
						}
						if _, err := mm.SendEmail(
							"plus users report",
							msg,
							"text/html",
							*recipientName,
							*emailRecipient,
						); err != nil {
							fmt.Println("failed to send email report", err.Error())
							os.Exit(1)
						}
					}
				},
			},
		},
	},
	"upload": {
		Blurb:         "Upload based metrics",
		Description:   "Allows for gathering of upload based metrics (number of uploads, type, etc...)",
		ChildRequired: true,
		Children: map[string]cmd.Cmd{
			"count": {
				Blurb:       "Upload count",
				Description: "Gets the total number of uploads",
				Action: func(cfg config.TemporalConfig, args map[string]string) {
					db, err := database.Initialize(&cfg, database.Options{
						SSLModeDisable: *dbNoSSL,
						RunMigrations:  *dbMigrate,
					})
					if err != nil {
						fmt.Println("failed to initialize database connection", err.Error())
						os.Exit(1)
					}
					ipfs, err := rtfs.NewManager(
						cfg.IPFS.APIConnection.Host+":"+cfg.IPFS.APIConnection.Port,
						"", 1*time.Minute,
					)
					if err != nil {
						fmt.Println("failed to open ipfs api connection", err.Error())
						os.Exit(1)
					}
					uf := upload.NewFarmer(db.DB, ipfs)
					num, err := uf.NumberOfUploads()
					if err != nil {
						fmt.Println("failed to get number of uploads", err.Error())
						os.Exit(1)
					}
					msg := fmt.Sprintf("there are %v total uploads", num)
					fmt.Println(msg)
					if *sendEmail {
						mm, err := mail.NewManager(&cfg, db.DB)
						if err != nil {
							fmt.Println("failed to initialize mail manager", err.Error())
							os.Exit(1)
						}
						if _, err := mm.SendEmail(
							"upload count report",
							msg,
							"text/html",
							*recipientName,
							*emailRecipient,
						); err != nil {
							fmt.Println("failed to send email report", err.Error())
							os.Exit(1)
						}
					}
				},
			},
			"size": {
				Blurb:       "Upload Size Average",
				Description: "Gets the average size of uploads",
				Action: func(cfg config.TemporalConfig, args map[string]string) {
					db, err := database.Initialize(&cfg, database.Options{
						SSLModeDisable: *dbNoSSL,
						RunMigrations:  *dbMigrate,
					})
					if err != nil {
						fmt.Println("failed to initialize database connection", err.Error())
						os.Exit(1)
					}
					ipfs, err := rtfs.NewManager(
						cfg.IPFS.APIConnection.Host+":"+cfg.IPFS.APIConnection.Port,
						"", 1*time.Minute,
					)
					if err != nil {
						fmt.Println("failed to open ipfs api connection", err.Error())
						os.Exit(1)
					}
					uf := upload.NewFarmer(db.DB, ipfs)
					size, err := uf.AverageUploadSize(*unique)
					if err != nil {
						fmt.Println("failed to get upload size average", err.Error())
						os.Exit(1)
					}
					var uniqueMessage string
					if *unique {
						uniqueMessage = "unique"
					} else {
						uniqueMessage = "non unique"
					}
					msg := fmt.Sprintf("the %s average size of uploads is %v gigabytes", uniqueMessage, size)
					fmt.Println(msg)
					if *sendEmail {
						mm, err := mail.NewManager(&cfg, db.DB)
						if err != nil {
							fmt.Println("failed to initialize mail manager", err.Error())
							os.Exit(1)
						}
						if resp, err := mm.SendEmail(
							"upload size report",
							msg,
							"text/html",
							*recipientName,
							*emailRecipient,
						); err != nil {
							fmt.Println("failed to send email report", err.Error())
							os.Exit(1)
						} else {
							fmt.Printf("resposne from sengrid\n+%v\n", resp)
						}
					}
				},
			},
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
