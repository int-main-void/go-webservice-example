package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

type Widget struct {
	ID   string
	Name string
}

type ExampleAppDb interface {
	GetWidgets() ([]Widget, error)
}

func GetAppDb() (ExampleAppDb, error) {
	if os.Getenv("DB_CLIENT") == "MOCK" {
		db, err := newMockAppDb()
		return db, err
	} else {
		db, err := newMysqlAppDb()
		return db, err
	}
}

//---------------------------------
// Mock Db Type

type MockAppDb struct {
	Widgets map[string]Widget
}

func newMockAppDb() (MockAppDb, error) {
	widgets := make(map[string]Widget)
	widgets["x"] = Widget{"1", "x"}
	widgets["y"] = Widget{"2", "y"}
	db := MockAppDb{widgets}
	return db, nil
}

func (db MockAppDb) GetWidgets() ([]Widget, error) {
	widgetset := make([]Widget, len(db.Widgets))
	i := 0
	for _, w := range db.Widgets {
		widgetset[i] = w
		i++
	}
	return widgetset, nil
}

//---------------------------------
// Real Mysql Db Type

type MysqlAppDb struct {
	db *sql.DB
}

func newMysqlAppDb() (MysqlAppDb, error) {
	dbHost := os.Getenv("MYSQL_HOST")
	dbPort := os.Getenv("MYSQL_PORT")
	dbUser := os.Getenv("MYSQL_USER")
	dbPw := os.Getenv("MYSQL_PW")
	dbName := os.Getenv("MYSQL_DB_NAME")

	//if dbHost == "" || dbPort

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true)", dbUser, dbPw, dbHost, dbPort, dbName)
	dbconn, err := sql.Open("mysql", connectionString)
	appdb := MysqlAppDb{db: dbconn}
	if err != nil {
		logrus.Error(err)
	}
	return appdb, err
}

func (db MysqlAppDb) GetWidgets() ([]Widget, error) {
	widgets := []Widget{}
	widgetsQueryString := "SELECT (id, name) FROM widgets"
	rows, err := db.db.Query(widgetsQueryString)
	defer rows.Close()
	if err != nil {
		logrus.Error(err)
		return widgets, err
	}
	for rows.Next() {
		var id string
		var name string
		rows.Scan(&id, &name)
		widgets = append(widgets, Widget{ID: id, Name: name})
	}
	return widgets, nil
}
