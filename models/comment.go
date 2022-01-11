package models

import (
	"log"
	"reflect"
)

type Comment struct {
	Id       int
	ParentId int
	ToUserId int
	UserId   int
}

func (co *Comment) Get(id int) (comment Comment, err error) {
	stout, err := Db.Prepare("select id,parent_id,user_id,to_user_id from comment where id=?  ")
	if err != nil {
		log.Printf("%+v\n", err)
	}
	defer stout.Close()
	err = stout.QueryRow(id).Scan(&comment.Id, &comment.ParentId, &comment.UserId, &comment.ToUserId)
	return
}
func (co Comment) IsEmpty() bool {
	return reflect.DeepEqual(co, Comment{})
}

func (co *Comment) Cos(id int) (comments []Comment, err error) {

	rows, err := Db.Query("select id,parent_id from comment where parent_id=? and is_delete=0 order by  digg_count desc, create_time desc limit 2", id)
	if err != nil {
		return
	}
	for rows.Next() {
		comment := Comment{}
		if err = rows.Scan(&comment.Id, &comment.ParentId); err != nil {
			return
		}
		comments = append(comments, comment)
	}
	rows.Close()
	return
}

func (co *Comment) Count(id int) (count int, err error) {
	stout, err := Db.Prepare("select count(id) as c from comment where parent_id=? and is_delete=0 ")
	if err != nil {
		log.Printf("%+v\n", err)
	}
	defer stout.Close()
	err = stout.QueryRow(id).Scan(&count)
	return
}
func (co *Comment) Update(str string, count int) (err error) {
	statement := "update comment set sub_ids = ?,sub_count=? where id = ? and is_delete = 0"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(str, count, co.ParentId)
	return
}
