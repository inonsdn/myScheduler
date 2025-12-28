package localdb

import (
	"database/sql"
	"fmt"
	"reflect"
	"scheduler/internal/config"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type DbConnection interface {
	execute(string, []any) (any, error)
	query(SqlResult, string, ...any) error
}

type LocalDb struct {
	opts       *config.Options
	driver     string
	dataSource string
	db         *sql.DB
}

type SqlResult interface {
	update(rows *sql.Rows) error
}

func NewLocalDb(opts *config.Options) (*LocalDb, error) {
	dsn := opts.DataSourceConfig()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Cannot open sql")
		return nil, err
	}
	return &LocalDb{
		opts:       opts,
		driver:     "mysql",
		dataSource: dsn,
		db:         db,
	}, nil
}

func (s *LocalDb) InitializeSchema() error {
	db, err := sql.Open(s.driver, s.dataSource)
	if err != nil {
		fmt.Println("Cannot open sql")
		return err
	}
	defer db.Close()

	for _, statement := range schemaStrList {
		_, err := s.db.Exec(statement)
		if err != nil {
			fmt.Println("Found error when init", err)
			break
		}
	}

	return nil
}

func (s *LocalDb) execute(statement string, args []any) (any, error) {

	stmt, err := s.db.Prepare(statement)
	if err != nil {
		return nil, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("Cannot execute statement.")
		return nil, err
	}

	return result.RowsAffected()
}

func (s *LocalDb) query(dest SqlResult, statement string, args ...any) error {

	rows, err := s.db.Query(statement, args...)
	defer rows.Close()

	if err != nil {
		fmt.Println("Found error when query", err)
		return nil
	}

	if !rows.Next() {
		// no rows returned
		if err := rows.Err(); err != nil {
			fmt.Println("rows error:", err)
			return err
		}
		// no error, just no data
		return sql.ErrNoRows
	}

	// return dest.update(rows)
	return update(dest, rows)
}

type TableBase interface{}

func GetAllColumnWithValue(t any) map[string]any {
	structType := reflect.TypeOf(t)
	structValue := reflect.ValueOf(t)

	// If pointer → get the underlying element
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
		structValue = structValue.Elem()
	}

	colToVal := make(map[string]any)
	fmt.Println("===== structType.NumField()", structType.NumField())

	// loop over field num
	for i := 0; i < structType.NumField(); i++ {

		// get field from index
		field := structType.Field(i)

		col := field.Tag.Get("sql")
		fmt.Println("===== col", col)
		if col != "" {
			colToVal[col] = structValue.Field(i).Interface()
		}
	}
	return colToVal
}

type SchedulerJob struct {
	TableBase
	Id      string `sql:"id"`
	Name    string `sql:"name"`
	JobType int    `sql:"jobType"`
	Year    int    `sql:"year"`
	Month   int    `sql:"month"`
	Day     int    `sql:"day"`
	Hour    int    `sql:"hour"`
	Minute  int    `sql:"minute"`
}

func NewReminderJob(name string, year int, month int, day int, hour int, minute int) *SchedulerJob {
	return &SchedulerJob{
		Id:      uuid.NewString(),
		Name:    name,
		JobType: 0,
		Year:    year,
		Month:   month,
		Day:     day,
		Hour:    hour,
		Minute:  minute,
	}
}
func update(dest SqlResult, rows *sql.Rows) error {
	cols, _ := rows.Columns()

	vals := make([]any, len(cols))
	valPtrs := make([]any, len(cols))

	// loop over array of value to map pointer array point to each value addr
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	if err := rows.Scan(valPtrs...); err != nil {
		fmt.Println("Found error when get from rows", err)
		return err
	}

	// reflect on destination
	rv := reflect.ValueOf(dest).Elem()
	rt := rv.Type()

	// map sql tag -> field index
	tagToField := map[string]int{}
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("sql")
		if tag != "" {
			tagToField[tag] = i
		}
	}

	// assign values by column name
	for i, colName := range cols {
		if fieldIndex, ok := tagToField[colName]; ok {
			rv.Field(fieldIndex).Set(reflect.ValueOf(vals[i]))
		}
	}

	return nil
}

// implement method update of interface sql result
func (s *SchedulerJob) update(rows *sql.Rows) error {

	err := rows.Scan(&s.Id, &s.Name)
	if err != nil {
		fmt.Println("Found error when get from rows", err)
		return err
	}
	return nil
}

type SchedulerJobTable struct {
	dbCon DbConnection
	name  string
}

func NewSchedulerJobTable(dbCon DbConnection) *SchedulerJobTable {
	return &SchedulerJobTable{
		dbCon: dbCon,
		name:  "schedulerJob",
	}
}

func (s *SchedulerJobTable) SetDbCon(dbCon DbConnection) {
	s.dbCon = dbCon
}

func (s *SchedulerJobTable) CreateSchedulerJob(sj *SchedulerJob) error {
	columnToValue := GetAllColumnWithValue(sj)
	var columns []string
	var placeholders []string
	var values []any
	for col, val := range columnToValue {
		columns = append(columns, col)
		values = append(values, val)
		placeholders = append(placeholders, "?")
	}

	columnsStr := strings.Join(columns, ",")
	params := strings.Join(placeholders, ",")

	statement := fmt.Sprintf(
		"INSERT INTO %s ( %s ) VALUES ( %s )", s.name, columnsStr, params)
	fmt.Println("Statement", statement)
	fmt.Println("Value", values)
	rowAffected, err := s.dbCon.execute(statement, values)

	fmt.Println("Executed, effect row", rowAffected)

	return err
}

func (s *SchedulerJobTable) GetAllJob() (*SchedulerJob, error) {
	schedulerJob := SchedulerJob{}
	statement := fmt.Sprintf(
		"SELECT id, name, jobType, year, month, day, hour, minute FROM %s",
		s.name,
	)
	err := s.dbCon.query(&schedulerJob, statement)
	return &schedulerJob, err
}
