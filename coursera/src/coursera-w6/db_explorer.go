package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type contextKey string

type query struct {
	table      string
	id         *int
	primaryKey string
	data       map[string]interface{}
	limit      int
	offset     int
}

type tableMetaData struct {
	Type       string
	Collation  sql.NullString
	Null       string
	Key        string
	Default    sql.NullBool
	Extra      sql.NullString
	Privileges sql.NullString
	Comment    string
}

func newQuery(table string) *query {
	return &query{
		table:  table,
		limit:  5,
		offset: 0,
	}
}

type dbConn struct {
	conn *sql.DB
}

// NewDbExplorer возвращает http handler для обработки запросов
func NewDbExplorer(conn *sql.DB) (http.Handler, error) {
	dbConn := dbConn{
		conn: conn,
	}

	return extractRequestParametersMiddleware(dbConn), nil
}

func getContextQueryKey(key contextKey) contextKey {
	var queryKey contextKey
	queryKey = key
	return queryKey
}

func extractRequestParametersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := regexp.MustCompile("^\\/(\\w+)(\\/(\\d+)?)?$")
		urlParams := pattern.FindStringSubmatch(r.URL.Path)

		countMatches := len(urlParams)

		if countMatches == 0 {
			next.ServeHTTP(w, r)
			return
		}

		tableName := urlParams[1]
		query := newQuery(tableName)

		rowID, err := strconv.Atoi(urlParams[3])
		if rowID > 0 && err == nil {
			query.id = &rowID
		}

		queryParams := r.URL.Query()
		query.limit = getLimitOffset("limit", queryParams, 5)
		query.offset = getLimitOffset("offset", queryParams, 0)

		if r.Method == http.MethodPut || r.Method == http.MethodPost {
			err := json.NewDecoder(r.Body).Decode(&query.data)
			if err != nil {
				panic(err)
			}
		}

		ctx := context.WithValue(r.Context(), getContextQueryKey("query"), query)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getLimitOffset(parameterName string, parameters url.Values, defaultValue int) int {
	value, ok := parameters[parameterName]
	parameter := defaultValue
	if ok {
		intValue, err := strconv.Atoi(value[0])
		if err == nil {
			parameter = intValue
		}
	}

	return parameter
}

func (dbConn dbConn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tableList := getTableList(dbConn.conn)

	if r.URL.Path == "/" && r.Method == http.MethodGet {
		writeResponse(w, tableList, "tables")
		return
	}

	query, ok := r.Context().Value(getContextQueryKey("query")).(*query)
	if !ok {
		panic(errors.New("Invalid URL, can't be parsed"))
	}

	err := checkTableNameIsCorrect(query.table, tableList)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, err)
		return
	}

	tmd := getTableMetaData(dbConn.conn, query.table)
	query.primaryKey = getPrimaryKeyField(tmd)

	if r.Method == http.MethodGet {
		rows := getTableRows(dbConn.conn, *query)

		if query.id == nil {
			writeResponse(w, rows, "records")
			return
		}

		if len(rows) == 0 {
			writeErrorResponse(w, http.StatusNotFound, errors.New("record not found"))
			return
		}
		writeResponse(w, rows[0], "record")
		return
	}

	if r.Method == http.MethodPut {
		result, err := insertTableRow(dbConn.conn, *query, tmd)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		writeResponse(w, result[0], "")
		return
	}

	if r.Method == http.MethodPost {
		result, err := updateTableRow(dbConn.conn, *query, tmd)
		if err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		writeResponse(w, result[0], "")
		return
	}

	if r.Method == http.MethodDelete {
		result := deleteTableRow(dbConn.conn, *query)
		writeResponse(w, result[0], "")
		return
	}
}

func getTableList(conn *sql.DB) []string {
	q := "SHOW TABLES;"
	rows, err := conn.Query(q)

	if err != nil {
		panic(err)
	}

	r := make([]string, 0)
	for rows.Next() {
		var col interface{}
		err = rows.Scan(&col)
		if err != nil {
			panic(err)
		}

		if v, ok := col.(string); ok {
			r = append(r, v)
		}

		if v, ok := col.([]byte); ok {
			r = append(r, string(v))
		}
	}

	rows.Close()

	return r
}

func getTableMetaData(conn *sql.DB, tableName string) map[string]tableMetaData {
	q := "SHOW FULL COLUMNS FROM " + tableName + ";"
	rows, err := conn.Query(q)
	if err != nil {
		panic(err)
	}

	result := make(map[string]tableMetaData, 0)
	for rows.Next() {
		var field sql.NullString
		tmd := tableMetaData{}

		err := rows.Scan(&field, &tmd.Type, &tmd.Collation, &tmd.Null, &tmd.Key, &tmd.Default, &tmd.Extra, &tmd.Privileges, &tmd.Comment)
		if err != nil {
			panic(err)
		}

		if field.Valid {
			result[field.String] = tmd
		}
	}

	rows.Close()
	return result
}

func getTableRows(conn *sql.DB, query query) []interface{} {
	q := "SELECT * FROM " + query.table
	if query.id != nil {
		q += " WHERE " + query.primaryKey + " = ?"
	}

	q += " LIMIT " + fmt.Sprint(query.limit) + " OFFSET " + fmt.Sprint(query.offset)

	stmt, err := conn.Prepare(q)
	// TODO обработать ошибку
	if err != nil {
		panic(err)
	}

	var rows *sql.Rows
	if query.id != nil {
		rows, err = stmt.Query(*query.id)
	} else {
		rows, err = stmt.Query()
	}
	// TODO обработать ошибку
	if err != nil {
		panic(err)
	}

	columns, err := rows.Columns()
	// TODO обработать ошибку
	if err != nil {
		panic(err)
	}

	columsCount := len(columns)

	result := make([]interface{}, 0)
	for rows.Next() {
		r := make([]interface{}, 0)
		for i := 0; i < columsCount; i++ {
			var in interface{}
			r = append(r, &in)
		}
		err := rows.Scan(r...)
		// TODO обработать ошибку
		if err != nil {
			panic(err)
		}

		m := make(map[string]interface{}, 0)
		for i := 0; i < columsCount; i++ {
			m[columns[i]] = r[i]
		}

		result = append(result, m)
	}

	rows.Close()

	return result
}

func insertTableRow(conn *sql.DB, query query, tmd map[string]tableMetaData) ([]interface{}, error) {
	var fieldList []string
	var fieldValues []interface{}
	var fieldParams []string

	for k, v := range tmd {
		if k == query.primaryKey && v.Extra.String == "auto_increment" {
			continue
		}

		value, ok := query.data[k]
		if !ok && v.Null == "YES" {
			continue
		}

		if !ok && !v.Default.Bool {
			fieldList = append(fieldList, k)
			fieldValues = append(fieldValues, "")
			fieldParams = append(fieldParams, "?")
			continue
		}

		var err error
		err = checkInputIntParameter(v, value, k)
		if err != nil {
			return make([]interface{}, 0), err
		}

		err = checkInputStringParameter(v, value, k)
		if err != nil {
			return make([]interface{}, 0), err
		}

		fieldList = append(fieldList, k)
		fieldValues = append(fieldValues, value)
		fieldParams = append(fieldParams, "?")
	}

	fields := strings.Join(fieldList, ",")
	fields = " (" + fields + ") "

	params := strings.Join(fieldParams, ",")
	params = " (" + params + ") "

	q := "INSERT INTO " + query.table + fields + "VALUES" + params + ";"

	stmt, err := conn.Prepare(q)
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	r, err := stmt.Exec(fieldValues...)
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	id, err := r.LastInsertId()
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{}, 0)
	m[query.primaryKey] = id

	result := make([]interface{}, 0)
	result = append(result, m)

	return result, nil
}

func checkInputIntParameter(fmd tableMetaData, v interface{}, k string) error {
	if fmd.Type == "int" {
		if v == nil && fmd.Null != "YES" {
			return errors.New("field " + k + " have invalid type")
		}
		if v != nil && reflect.TypeOf(v).Kind() != reflect.Int {
			return errors.New("field " + k + " have invalid type")
		}
	}

	return nil
}

func checkInputStringParameter(fmd tableMetaData, v interface{}, k string) error {
	if fmd.Type == "text" || strings.Contains(fmd.Type, "varchar") {
		if v == nil && fmd.Null != "YES" {
			return errors.New("field " + k + " have invalid type")
		}
		if v != nil && reflect.TypeOf(v).Kind() != reflect.String {
			return errors.New("field " + k + " have invalid type")
		}
	}

	return nil
}

func updateTableRow(conn *sql.DB, query query, tmd map[string]tableMetaData) ([]interface{}, error) {
	var fieldList []string
	var fieldValues []interface{}
	var fieldParams []string

	for k, v := range query.data {
		if k == query.primaryKey {
			return make([]interface{}, 0), errors.New("field " + query.primaryKey + " have invalid type")
		}

		fieldMetaData := tmd[k]

		var err error
		err = checkInputIntParameter(fieldMetaData, v, k)
		if err != nil {
			return make([]interface{}, 0), err
		}

		err = checkInputStringParameter(fieldMetaData, v, k)
		if err != nil {
			return make([]interface{}, 0), err
		}

		field := k + " = ?"
		fieldList = append(fieldList, field)
		fieldValues = append(fieldValues, v)
		fieldParams = append(fieldParams, "?")
	}

	fields := strings.Join(fieldList, ",")

	q := "UPDATE " + query.table + " SET " + fields

	if query.id != nil {
		//Добавляем еще один параметр для WHERE
		fieldValues = append(fieldValues, query.id)
		q += " WHERE " + query.primaryKey + " = ?;"
	}

	stmt, err := conn.Prepare(q)
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	r, err := stmt.Exec(fieldValues...)
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	affectedCount, err := r.RowsAffected()
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{}, 0)
	m["updated"] = affectedCount

	result := make([]interface{}, 0)
	result = append(result, m)

	return result, nil
}

func deleteTableRow(conn *sql.DB, query query) []interface{} {
	q := "DELETE FROM " + query.table

	if query.id != nil {
		q += " WHERE " + query.primaryKey + " = ?;"
	}

	stmt, err := conn.Prepare(q)
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	r, err := stmt.Exec(query.id)
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	affectedCount, err := r.RowsAffected()
	// TODO ошибка
	if err != nil {
		panic(err)
	}

	m := make(map[string]interface{}, 0)
	m["deleted"] = affectedCount

	result := make([]interface{}, 0)
	result = append(result, m)

	return result
}

func getPrimaryKeyField(tmd map[string]tableMetaData) string {
	for k, v := range tmd {
		if v.Key == "PRI" {
			return k
		}
	}

	return ""
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	result := make(map[string]interface{}, 0)
	result["error"] = err.Error()

	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(status)
	io.WriteString(w, string(jsonResult))
}

func writeResponse(w http.ResponseWriter, input interface{}, prefix string) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, makeResponse(input, prefix))
}

// MakeResponse преобразует результат из вида слайс интерфейсов к виду JSON
func makeResponse(input interface{}, prefix string) string {
	inputValue := reflect.ValueOf(input)
	unpack(inputValue)

	var jsonResult []byte
	if inputValue.Kind() == reflect.Slice {
		jsonResult = makeSliceResponse(inputValue, prefix)
	}

	if inputValue.Kind() == reflect.Map {
		if prefix == "" {
			jsonResult = makeModifyResponse(inputValue)
		} else {
			jsonResult = makeSingleResponse(inputValue, prefix)
		}
	}

	return string(jsonResult)
}

func makeSliceResponse(value reflect.Value, prefix string) []byte {
	result := make(map[string]map[string][]interface{}, 0)
	result["response"] = make(map[string][]interface{}, 0)
	result["response"][prefix] = make([]interface{}, 0)

	for i := 0; i < value.Len(); i++ {
		result["response"][prefix] = append(result["response"][prefix], value.Index(i).Interface())
	}

	jsonResult, err := json.Marshal(result)

	if err != nil {
		panic(err)
	}

	return jsonResult
}

func makeSingleResponse(value reflect.Value, prefix string) []byte {
	result := make(map[string]map[string]interface{}, 0)
	result["response"] = make(map[string]interface{}, 0)

	line := make(map[string]interface{}, 0)
	iter := value.MapRange()
	for iter.Next() {
		line[iter.Key().String()] = iter.Value().Interface()
	}
	result["response"][prefix] = line

	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return jsonResult
}

func makeModifyResponse(value reflect.Value) []byte {
	result := make(map[string]interface{}, 0)

	result["response"] = value.Interface()
	jsonResult, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}

	return jsonResult
}

// Распаковывает и заменяет значения, которые лежат в интерефейсах для сложных типов (slice|map)
func unpack(value reflect.Value) {
	if value.Kind() == reflect.Slice {
		for i := 0; i < value.Len(); i++ {
			unpackValue(value.Index(i))
		}
	}

	if value.Kind() == reflect.Map {
		iter := value.MapRange()
		for iter.Next() {
			unpackValue(iter.Value())
		}
	}
}

// Распаковывает и заменяет значения, которые лежат в интерефейсах
func unpackValue(value reflect.Value) {
	element := value.Interface()
	sliceElement := reflect.ValueOf(element)
	if sliceElement.Kind() == reflect.Slice {
		unpack(sliceElement)
	}

	if sliceElement.Kind() == reflect.Map {
		unpack(sliceElement)
	}

	elementValue := reflect.ValueOf(element)

	if elementValue.Kind() == reflect.Ptr {
		value := elementValue.Elem()
		if v, ok := value.Interface().([]byte); ok {
			value.Set(reflect.ValueOf(string(v)))
		}
	}
}

func checkTableNameIsCorrect(tableName string, tableList []string) error {
	for _, table := range tableList {
		if tableName == table {
			return nil
		}
	}

	return errors.New("unknown table")
}
