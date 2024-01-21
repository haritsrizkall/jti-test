package pkg

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func (m *MySQL) Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", m.Username, m.Password, m.Host, m.Port, m.Database))
	if err != nil {
		return nil, err
	}

	return db, nil
}
