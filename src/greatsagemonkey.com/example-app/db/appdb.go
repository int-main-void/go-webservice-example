package db

import "errors"

type Widget struct {
	ID   string
	Name string
}

type ExampleAppDb interface {
	GetWidgets() ([]Widget, error)
}

func GetAppDb() ExampleAppDb {
	// TODO - switch impl based on env, don't reload, etc.
	return NewMockAppDb()
}

//---------------------------------
// Mock Db Type

type MockAppDb struct {
	Widgets map[string]Widget
}

func NewMockAppDb() MockAppDb {
	widgets := make(map[string]Widget)
	widgets["x"] = Widget{"1", "x"}
	widgets["y"] = Widget{"2", "y"}
	db := MockAppDb{widgets}
	return db
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
}

func NewMysqlAppDb() {

}

func (db *MysqlAppDb) GetWidgets() ([]Widget, error) {
	return nil, errors.New("not implemented")
}
