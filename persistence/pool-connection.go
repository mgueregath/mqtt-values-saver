package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/kisielk/sqlstruct"
	"log"
	"reflect"
)

type PoolConnection struct {
	connection *sql.DB
}

func NewPoolConnection(connection *sql.DB) *PoolConnection {
	return &PoolConnection{connection}
}

func (poolConnection *PoolConnection) Query(queryString string, dest interface{}) error {
	rows, err := poolConnection.connection.Query(queryString)
	if err != nil {
		return err
	}
	//defer rows.Close()

	slice := reflect.ValueOf(dest)
	if slice.Kind() != reflect.Ptr {
		return errors.New("dest must be pointer")
	}

	slice = slice.Elem()
	if slice.Kind() != reflect.Slice {
		return errors.New("dest must be pointer to struct")
	}

	elementp := reflect.New(slice.Type().Elem())

	for rows.Next() {

		err := sqlstruct.Scan(elementp.Interface(), rows)
		if err != nil {
			log.Fatalln(err)
		}
		slice.Set(reflect.Append(slice, elementp.Elem()))
	}

	return nil
}

func (poolConnection *PoolConnection) Save(queryString string, values []interface{}) error {
	insert, err := poolConnection.connection.Query(fmt.Sprintf(queryString, values...))
	if err != nil {
		return err
	}
	defer insert.Close()

	return nil
}
