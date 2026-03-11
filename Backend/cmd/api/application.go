package main

import "fmt"

type Application struct {
	config Config
}

type Config struct {
	dbConfig DBConfig
}

type DBConfig struct {
	username string
	password string
	dbName   string
	dbHost   string
	dbAddr   string
}

func (app *Application) Mount() {
	fmt.Println(app.config)
}
