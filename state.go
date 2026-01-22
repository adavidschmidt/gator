package main

import (
	"github.com/adavidschmidt/blogaggregator/internal/config"
	"github.com/adavidschmidt/blogaggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}
