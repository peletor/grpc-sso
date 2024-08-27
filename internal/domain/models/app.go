package models

type App struct {
	ID     int
	Name   string
	Secret string
}

const EmptyAppID = 0
