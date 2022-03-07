// Package database provides some functions to interface with the posstgresql database
package database

import (
	"github.com/soumitradev/Dwitter/backend/common"
	"github.com/soumitradev/Dwitter/backend/prisma/db"
)

func init() {
	ConnectDB()
}

// Connect to the database using prisma
func ConnectDB() {
	common.Client = db.NewClient()
	if err := common.Client.Prisma.Connect(); err != nil {
		panic(err)
	}
}

// Disconnect from DB
func DisconnectDB() {
	if err := common.Client.Prisma.Disconnect(); err != nil {
		panic(err)
	}
}
