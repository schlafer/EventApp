package main

import (
	"database/sql"
	"log"

	_ "github.com/joho/godotenv/autoload" // Automatically loads environment variables
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/schlafer/EventApp/docs"
	"github.com/schlafer/EventApp/internal/database"
	"github.com/schlafer/EventApp/internal/env"
)

// @title EventApp API Documentation
// @version 1.0
// @description A rest API in Go using Gin framework.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// Apply the security definition to your endpoints
// @security BearerAuth

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {

	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	models := database.NewModels(db)

	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "123secret"),
		models:    models,
	}

	if err := serve(app); err != nil {
		log.Fatal(err)
	}
}

/*
Here we load environment variables, initialize the database connection,
create an application struct and start the server using the serve function.
The application struct will be used to pass the dependencies around
without having global variables.
We then start the server using the serve function.
*/
