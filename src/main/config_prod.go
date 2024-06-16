//go:build prod
// +build prod

package main

import (
	"os"
)

var (
	DB_USER      = getEnv("MYSQL_DB_USERNAME", "backend_db_client")
	DB_PASS      = getEnv("MYSQL_DB_PASSWORD", "proneraggedyplanetgallows")
	DB_NAME      = getEnv("MYSQL_DATABASE", "cookbook")
	DB_HOST      = getEnv("MYSQL_HOST", "cookbookDB")
	BACKEND_PORT = getEnv("BACKEND_PORT", "8085")
	DB_PORT      = getEnv("MYSQL_DB_PORT", "3307")
)

// getEnv gets the value of the given environment variable or returns the provided default value if the environment variable is not set.
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
