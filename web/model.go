package web

import (
	"gonote/internal/app"
	"gonote/pkg/utils"
)

type Note struct {
	Name string
	Num  int
	Uuid string
}

func insert(name string, num int) error {
	uuid := utils.Rand()
	_, err := app.Db.Exec("insert into QRTZ_NOTE(`name`, `num`, `uuid`) values(?, ?, ?)", name, num, uuid.Hex())
	return err
}

func update(name string, num int, uuid string) error {
	_, err := app.Db.Exec("update QRTZ_NOTE set `name`=?, `num` = ? where uuid = ?", name, num, uuid)
	return err
}

func delete(uuid string) error {
	_, err := app.Db.Exec("delete from QRTZ_NOTE where uuid = ?", uuid)
	return err
}

func get() []Note {
	rows, _ := app.Db.Query("select `name`, `num`, `uuid` from QRTZ_NOTE")
	defer rows.Close()
	noteList := make([]Note, 0)
	for rows.Next() {
		var name string
		var num int
		var uuid string
		rows.Scan(&name, &num, &uuid)
		noteList = append(noteList, Note{Name: name, Num: num, Uuid: uuid})
	}
	return noteList
}
