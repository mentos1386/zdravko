package models

type Healthcheck struct {
	ID               uint `gorm:"primary_key"`
	Name             string
	Status           string // UP, DOWN
	UptimePercentage float64
}
