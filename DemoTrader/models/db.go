package models

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbpath = "./db.sqlite3"
	dbtype = "sqlite3"
)

func connect() *sqlx.DB {
	db, err := sqlx.Connect(dbtype, dbpath)
	if err != nil {
		fmt.Println("---- connect error ----")
		log.Fatal(err)
	}
	return db
}

func insertQuery(tableName string, entryData map[string]interface{}) string {
	var stringColumns string
	var stringValues string
	var stringColumnsList []string
	var stringValuesList []string

	for k := range entryData {
		stringColumnsList = append(stringColumnsList, k)

		if fmt.Sprintf("%T", entryData[k]) == "int" {
			stringValuesList = append(stringValuesList, fmt.Sprintf("%d", entryData[k]))
		} else {
			stringValuesList = append(stringValuesList, fmt.Sprintf("'%s'", entryData[k]))
		}

	}

	stringColumns = strings.Join(stringColumnsList, ", ")
	stringValues = strings.Join(stringValuesList, ", ")

	query := "INSERT INTO %s (%s) VALUES (%s)"
	query = fmt.Sprintf(query, tableName, stringColumns, stringValues)
	return query
}

func selectQuery(readFields []string, tableName string, conditionData map[string]interface{}) string {
	var stringValues string
	var stringValuesList []string

	for k := range conditionData {
		stringValuesList = append(stringValuesList, fmt.Sprintf("'%s'", conditionData[k]))

	}
	stringReadFields := "*"
	if len(readFields) > 0 && readFields[0] != "*" {
		stringReadFields = strings.Join(readFields, ", ")
	}

	stringValues = strings.Join(stringValuesList, " AND ")
	query := "SELECT %s FROM %s"
	if len(stringValuesList) > 0 {
		query += " WHERE %s"
		query = fmt.Sprintf(query, stringReadFields, tableName, stringValues)
	} else {
		query = fmt.Sprintf(query, stringReadFields, tableName)
	}

	return query
}

func insertDB(query string) {
	db := connect()
	db.Exec(query)
}

func selectDB(query string) *sql.Rows {
	db := connect()
	results, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
	return results
}
