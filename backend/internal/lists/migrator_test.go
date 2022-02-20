package lists

import (
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	databaseName = "test_rwlist_engine"
)

type databaseMigrator struct {
	t               *testing.T
	db              *gorm.DB
	connstr         string
	currentDatabase string
}

func newDatabaseMigrator(t *testing.T, connstr string) *databaseMigrator {
	var db *gorm.DB
	var err error

	for retry := 0; retry < 10; retry++ {
		db, err = gorm.Open(postgres.Open(connstr), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		t.Fatal("failed to connect to postgres", err)
	}
	db = db.Debug()

	currentDatabase := db.Migrator().CurrentDatabase()
	t.Log("connected to database", currentDatabase)

	return &databaseMigrator{
		t:               t,
		db:              db,
		connstr:         connstr,
		currentDatabase: currentDatabase,
	}
}

func (m *databaseMigrator) makeFreshDatabase() {
	err := m.createDatabase()
	if err != nil {
		// try to drop and retry
		err = m.dropIfExists()
		if err != nil {
			m.t.Fatal("failed to drop existing database", err)
		}
		err = m.createDatabase()
		if err != nil {
			m.t.Fatal("failed to create database", err)
		}
	}
}

func (m *databaseMigrator) createDatabase() error {
	return m.db.Exec("CREATE DATABASE " + databaseName).Error
}

func (m *databaseMigrator) dropIfExists() error {
	return m.db.Exec("DROP DATABASE IF EXISTS " + databaseName).Error
}

func (m *databaseMigrator) close() {
	err := m.dropIfExists()
	if err != nil {
		m.t.Error("failed to drop database after tests", err)
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		m.t.Fatal("failed to get sql db", err)
	}
	err = sqlDB.Close()
	if err != nil {
		m.t.Fatal("failed to close sql db", err)
	}
}

func (m *databaseMigrator) connectNew() *gorm.DB {
	connstr := strings.Replace(m.connstr, "dbname="+m.currentDatabase, "dbname="+databaseName, 1)

	db, err := gorm.Open(postgres.Open(connstr), &gorm.Config{})
	if err != nil {
		m.t.Fatal("failed to connect to postgres", err)
	}
	db = db.Debug()

	m.t.Log("connected to database", db.Migrator().CurrentDatabase())
	return db
}
