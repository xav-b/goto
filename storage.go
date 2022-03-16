package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const DB_DRIVER string = "sqlite3"

type Storage struct {
	db *sql.DB
}

func NewStorage(dbPath string, reset bool) (*Storage, error) {
	if reset {
		log.Println("reseting database")
		os.Remove(dbPath)
	}

	log.Printf("opening db %s\n", dbPath)
	db, err := sql.Open(DB_DRIVER, dbPath)

	return &Storage{db}, err
}

func (s *Storage) Init() error {
	sql_table := `
	CREATE TABLE IF NOT EXISTS service (
		id INTEGER primary KEY AUTOINCREMENT,

		name TEXT,           -- Github
		link TEXT NOT NULL,  -- https://github.com/xav-b/goto
		alias TEXT NOT NULL, -- git/xav-b/goto
		description TEXT,    -- a simple example
		tags TEXT,           -- example,vcs

		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	log.Printf("initialising service table if not exists\n")
	if _, err := s.db.Exec(sql_table); err != nil {
		return err
	}

	sql_table = `
	CREATE TABLE IF NOT EXISTS log (
		id INTEGER primary KEY AUTOINCREMENT,

		-- NOTE could be a foreign key
		link TEXT NOT NULL,
		alias TEST NOT NULL,
		user TEXT,

		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	log.Printf("initialising slog table if not exists\n")
	if _, err := s.db.Exec(sql_table); err != nil {
		return err
	}

	log.Println("successfully initialised DB")
	return nil
}

func (s *Storage) byAlias(alias string) *Service {
	sqlSelect := `
		SELECT link, alias
		FROM service
		WHERE alias = ?
	`

	log.Printf("fetching service (alias=%s)\n", alias)
	item := &Service{}
	if err := s.db.QueryRow(sqlSelect, alias).Scan(&item.Link, &item.Alias); err != nil {
		// could also mean no results
		log.Printf("failed to scan service row: %v\n", err)
		return nil
	}

	return item
}

func (s *Storage) List(limit int) (results []*Service) {
	// NOTE: allow to limit to prefix? Like `Jira`
	sqlReadall := fmt.Sprintf(`
		SELECT id, link, alias, description, tags, created_at
		FROM service
		ORDER BY datetime(created_at) DESC
		LIMIT %d
	`, limit)

	log.Printf("fetching user services (limit=%d)\n", limit)
	rows, _ := s.db.Query(sqlReadall)
	defer rows.Close()

	for rows.Next() {
		var id int
		var tags string
		item := &Service{}
		if err := rows.Scan(&id, &item.Link, &item.Alias, &item.Description, &tags, &item.CreatedAt); err != nil {
			log.Fatalf("failed to scan row: %v\n", err)
		}
		item.Tags = strings.Split(tags, ",")

		log.Printf("found new item #%d: %v\n", id, item)
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		// TODO: should return err
		log.Fatalf("failed to fetch services: %v\n")
	}

	return results
}

func (s *Storage) AliasService(service *Service) error {
	// TODO: handle conflicts
	sqlAdd := `
		INSERT INTO service	(
			link, alias, description, tags
		) VALUES (?, ?, ?, ?)
	`

	stmt, _ := s.db.Prepare(sqlAdd)
	defer stmt.Close()

	_, err := stmt.Exec(service.Link, service.Alias, service.Description, strings.Join(service.Tags, ","))

	return err
}
