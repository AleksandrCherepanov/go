package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

const InvalidToken = "Invalid Token"
const ValidToken = "Valid Token"
const TimeOutToken = "Timeout Token"

type XMLRoot struct {
	XMLName xml.Name `xml:"root"`
	Rows    []XMLRow `xml:"row"`
}

type XMLRow struct {
	ID        int    `xml:"id"`
	Age       int    `xml:"age"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Name      string
}

func (r XMLRow) getName() string {
	return strings.Join([]string{r.FirstName, r.LastName}, " ")
}

func sortByName(rows []XMLRow, order int) []XMLRow {
	if order < 0 {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].getName() > rows[j].getName()
		})
	} else if order > 0 {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].getName() < rows[j].getName()
		})
	}

	return rows
}

func sortByAge(rows []XMLRow, order int) []XMLRow {
	if order < 0 {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].Age > rows[j].Age
		})
	} else if order > 0 {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].Age < rows[j].Age
		})
	}

	return rows
}

func sortByID(rows []XMLRow, order int) []XMLRow {
	if order < 0 {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].ID > rows[j].ID
		})
	} else if order > 0 {
		sort.Slice(rows, func(i, j int) bool {
			return rows[i].ID < rows[j].ID
		})
	}

	return rows
}

func sortSlice(rows []XMLRow, orderField string, order int) ([]XMLRow, error) {
	switch orderField {
	case "Name":
		return sortByName(rows, order), nil
	case "Id":
		return sortByID(rows, order), nil
	case "Age":
		return sortByAge(rows, order), nil
	case "":
		return sortByName(rows, order), nil
	default:
		e := errors.New("ErrorBadOrderField")
		return nil, e
	}
}

func searchQuery(rows []XMLRow, query string) []XMLRow {
	if query == "" {
		return rows
	}

	searchResult := make([]XMLRow, 0, len(rows))
	for _, row := range rows {
		lowerCaseName := strings.ToLower(row.getName())
		lowerCaseQuery := strings.ToLower(query)
		lowerCaseAbout := strings.ToLower(row.About)

		if strings.Contains(lowerCaseName, lowerCaseQuery) || strings.Contains(lowerCaseAbout, lowerCaseQuery) {
			searchResult = append(searchResult, row)
		}
	}

	return searchResult
}

func ServerSerach(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("AccessToken")
	if accessToken == InvalidToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if accessToken == TimeOutToken {
		time.Sleep(time.Second * 2)
		return
	}

	file, err := os.OpenFile("dataset.xml", os.O_RDONLY, os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var xmlData XMLRoot
	err = xml.Unmarshal(fileData, &xmlData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	rows := xmlData.Rows
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	query := r.URL.Query().Get("query")
	orderField := r.URL.Query().Get("order_field")
	orderBy, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if orderBy > 1 {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "order_by must be 1, -1, or 0")
		return
	}

	if orderBy < -1 {
		w.WriteHeader(http.StatusBadRequest)
		searchErrorResponse := SearchErrorResponse{
			Error: "order_by must be 1, -1, or 0",
		}
		errorResponse, _ := json.Marshal(searchErrorResponse)
		io.WriteString(w, string(errorResponse))
		return
	}

	rows = searchQuery(rows, query)
	rows, err = sortSlice(rows, orderField, orderBy)
	if err != nil {
		searchErrorResponse := SearchErrorResponse{
			Error: err.Error(),
		}
		errorResponse, err := json.Marshal(searchErrorResponse)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, string(errorResponse))
		return
	}

	if offset > len(rows) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if offset > 0 {
		rows = rows[offset:]
	}

	if limit > 0 {
		if limit < len(rows) {
			rows = rows[0:limit]
		}
	}

	result := make([]XMLRow, 0, len(rows))
	for _, row := range rows {
		resultRow := row
		resultRow.Name = row.getName()
		result = append(result, resultRow)
	}

	jsonResponse, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(rows) == 0 {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Nothing was found")
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResponse))
}

func makeTestRequest(limit, offset int, query, orderFiled string, orderBy int) SearchRequest {
	return SearchRequest{
		Limit:      limit,
		Offset:     offset,
		Query:      query,
		OrderField: orderFiled,
		OrderBy:    orderBy,
	}
}

func TestServerSearchUnauthorize(t *testing.T) {
	const expectedError = "Bad AccessToken"
	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: InvalidToken,
	}

	request := makeTestRequest(1, 0, "", "", 0)

	res, err := client.FindUsers(request)
	if res != nil {
		t.Errorf("Unexpected response. \nExpected: %v. \nActual: %v", nil, res)
	}
	if err.Error() != expectedError {
		t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", expectedError, err.Error())
	}

	ts.Close()
}

func TestServerSearchTimeout(t *testing.T) {
	const expectedError = "timeout for limit=11&offset=5&order_by=-1&order_field=Id&query=John"
	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: TimeOutToken,
	}

	request := makeTestRequest(10, 5, "John", "Id", -1)

	res, err := client.FindUsers(request)
	if res != nil {
		t.Errorf("Unexpected response. \nExpected: %v. \nActual: %v", nil, res)
	}
	if err.Error() != expectedError {
		t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", expectedError, err.Error())
	}

	ts.Close()
}

func TestServerSearchUnknownError(t *testing.T) {
	const expectedError = `unknown error Get ?limit=11&offset=5&order_by=-1&order_field=Id&query=John: unsupported protocol scheme ""`
	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         "",
		AccessToken: TimeOutToken,
	}

	request := makeTestRequest(10, 5, "John", "Id", -1)

	res, err := client.FindUsers(request)
	if res != nil {
		t.Errorf("Unexpected response. \nExpected: %v. \nActual: %v", nil, res)
	}
	if err.Error() != expectedError {
		t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", expectedError, err.Error())
	}

	ts.Close()
}

func TestServerSearchLimitOffsetIncorrect(t *testing.T) {
	type testCase struct {
		testCase string
		limit    int
		offset   int
		expected string
	}

	testCases := []testCase{
		{"Invalid Limit", -1, 5, "limit must be > 0"},
		{"Invalid Offset", 5, -1, "offset must be > 0"},
	}

	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	for _, test := range testCases {
		request := makeTestRequest(test.limit, test.offset, "John", "Id", -1)
		res, err := client.FindUsers(request)
		if res != nil {
			t.Errorf("Unexpected response. \nExpected: %v. \nActual: %v", nil, res)
		}

		if err.Error() != test.expected {
			t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", test.expected, err.Error())
		}
	}

	ts.Close()
}

func TestServerSearchInternalServerError(t *testing.T) {
	const expectedError = "SearchServer fatal error"
	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	request := makeTestRequest(5, 50, "", "", 0)

	res, err := client.FindUsers(request)
	if res != nil {
		t.Errorf("Unexpected response. \nExpected: %v. \nActual: %v", nil, res)
	}
	if err.Error() != expectedError {
		t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", expectedError, err.Error())
	}

	ts.Close()
}

func TestServerSearchBadRequest(t *testing.T) {
	type testCase struct {
		testCase   string
		orderField string
		orderBy    int
		expected   string
	}

	testCases := []testCase{
		{"Invalid Order Field", "Gender", -1, "OrderFeld Gender invalid"},
		{"Invalid Order By Incrorrect Error", "Id", 2, "cant unpack error json: invalid character 'o' looking for beginning of value"},
		{"Invalid Order By Correct Error", "Id", -2, "unknown bad request error: order_by must be 1, -1, or 0"},
	}

	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	for _, test := range testCases {
		request := makeTestRequest(5, 0, "John", test.orderField, test.orderBy)
		res, err := client.FindUsers(request)
		if res != nil {
			t.Errorf("Unexpected response. \nExpected: %v. \nActual: %v", nil, res)
		}

		if err.Error() != test.expected {
			t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", test.expected, err.Error())
		}
	}

	ts.Close()
}

func TestServerSearchResponseIsNotJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	request := makeTestRequest(25, 0, "Alexander", "", 0)
	res, err := client.FindUsers(request)
	if res != nil {
		t.Errorf("Unexpected response. Response should be nil")
	}

	if err.Error() != "cant unpack result json: invalid character 'N' looking for beginning of value" {
		t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", "cant unpack result json: invalid character 'N' looking for beginning of value", err.Error())
	}

	ts.Close()
}

func TestServerSearchSuccess(t *testing.T) {
	type testCase struct {
		testCase   string
		limit      int
		offset     int
		query      string
		orderField string
		orderBy    int
	}

	testCases := []testCase{
		{"Limit > 25", 30, 0, "", "", -1},
		{"Limit > than all data", 10, 30, "", "", 0},
		{"Search by name", 25, 0, "Guerr", "Id", -1},
		{"Search by about", 25, 0, "Lorem", "Age", 1},
	}

	ts := httptest.NewServer(http.HandlerFunc(ServerSerach))

	client := SearchClient{
		URL:         ts.URL,
		AccessToken: ValidToken,
	}

	for _, test := range testCases {
		request := makeTestRequest(test.limit, test.offset, test.query, test.orderField, test.orderBy)
		res, err := client.FindUsers(request)
		if res == nil {
			t.Errorf("Unexpected response. Response should not be nil")
		}

		if err != nil {
			t.Errorf("Unexpected error. \nExpected: %v. \nActual: %v", nil, err.Error())
		}
	}

	ts.Close()
}
