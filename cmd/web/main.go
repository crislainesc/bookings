package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/crislainesc/bookings/internal/config"
	"github.com/crislainesc/bookings/internal/driver"
	"github.com/crislainesc/bookings/internal/handlers"
	"github.com/crislainesc/bookings/internal/helpers"
	"github.com/crislainesc/bookings/internal/models"
	"github.com/crislainesc/bookings/internal/render"
	"github.com/joho/godotenv"
)

const (
	portNumber = ":8080"
)

var (
	app      config.AppConfig
	session  *scs.SessionManager
	infoLog  *log.Logger
	errorLog *log.Logger
)

// main is the main function
func main() {
	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()
	defer close(app.MailChan)

	listenForMail()

	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Starting application on port %s\n", portNumber)

	server := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

// loadEnv loads the environment variables from the .env file.
func loadEnv(envFile string) {
	err := godotenv.Load(dir(envFile))
	if err != nil {
		panic(fmt.Errorf("error loading .env file: %w", err))
	}
}

/*
dir returns the absolute path of the given environment file (envFile) in the Go module's
root directory. It searches for the 'go.mod' file from the current working directory upwards
and appends the envFile to the directory containing 'go.mod'.
It panics if it fails to find the 'go.mod' file.
*/
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}

func run() (*driver.Database, error) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(map[string]int{})

	loadEnv(".env")

	inProduction := os.Getenv("IN_PRODUCTION")
	useCache := os.Getenv("USE_CACHE")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbSSL := os.Getenv("DB_SSL")

	if dbName == "" || dbUser == "" {
		fmt.Println("Missing required flags")
		os.Exit(1)
	}

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// change this to true when in production
	app.InProduction, _ = strconv.ParseBool(inProduction)

	infoLog = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	log.Println("Connecting to database...")
	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		dbHost,
		dbPort,
		dbName,
		dbUser,
		dbPassword,
		dbSSL,
	)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	log.Println("Connected to database!")

	tcache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache", err.Error())
		return nil, err
	}

	app.TemplateCache = tcache
	app.UseCache, _ = strconv.ParseBool(useCache)

	repository := handlers.NewRepository(&app, db)
	helpers.NewHelpers(&app)
	handlers.NewHandlers(repository)

	render.NewRenderer(&app)
	return db, nil
}
