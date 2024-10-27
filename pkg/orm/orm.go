package orm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	ErrNoRows          = errors.New("no rows found")
	ErrInvalidOperator = errors.New("invalid operator")
	ErrInvalidValue    = errors.New("invalid value")
)

// DB represents the database connection
type Orm struct {
	*sql.DB
	mu       sync.RWMutex
	queryLog bool
	prepared map[string]*sql.Stmt
}

// Query represents a database query builder
type Query struct {
	table      string
	selections []string
	wheres     []whereClause
	orWheres   []whereClause
	joins      []joinClause
	limit      int
	offset     int
	orderBy    string
	orderDir   string
	groupBy    []string
	having     []havingClause
}

// Model represents a database model
type Model struct {
	db    *Orm
	query Query
	ctx   context.Context
}

type whereClause struct {
	column   string
	operator string
	value    interface{}
}

type joinClause struct {
	table     string
	joinType  string
	condition string
	args      []interface{}
}

type havingClause struct {
	condition string
	args      []interface{}
}

// Valid operators for where clauses
var validOperators = map[string]bool{
	"=":           true,
	"<>":          true,
	">":           true,
	"<":           true,
	">=":          true,
	"<=":          true,
	"LIKE":        true,
	"NOT LIKE":    true,
	"IN":          true,
	"NOT IN":      true,
	"IS NULL":     true,
	"IS NOT NULL": true,
}

// Config represents database configuration
type Config struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	QueryLog        bool
}

// New creates a new ORM instance with configuration
func New(db *sql.DB, config Config) *Orm {
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	return &Orm{
		DB:       db,
		queryLog: config.QueryLog,
		prepared: make(map[string]*sql.Stmt),
	}
}

// WithContext adds context to the model
func (m *Model) WithContext(ctx context.Context) *Model {
	m.ctx = ctx
	return m
}

// Table initializes a new query for the given table
func (db *Orm) Table(tableName string) *Model {
	return &Model{
		db:  db,
		ctx: context.Background(),
		query: Query{
			table:      tableName,
			selections: []string{fmt.Sprintf("%s.*", tableName)},
		},
	}
}

// Select adds columns to select
func (m *Model) Select(columns ...string) *Model {
	m.query.selections = sanitizeColumns(columns)
	return m
}

// Join methods with improved error handling and sanitization
func (m *Model) Join(table string, condition string, args ...interface{}) *Model {
	return m.addJoin("INNER JOIN", sanitizeTableName(table), condition, args...)
}

func (m *Model) LeftJoin(table string, condition string, args ...interface{}) *Model {
	return m.addJoin("LEFT JOIN", sanitizeTableName(table), condition, args...)
}

func (m *Model) RightJoin(table string, condition string, args ...interface{}) *Model {
	return m.addJoin("RIGHT JOIN", sanitizeTableName(table), condition, args...)
}

func (m *Model) CrossJoin(table string) *Model {
	return m.addJoin("CROSS JOIN", sanitizeTableName(table), "", nil)
}

func (m *Model) addJoin(joinType string, table string, condition string, args ...interface{}) *Model {
	m.query.joins = append(m.query.joins, joinClause{
		table:     table,
		joinType:  joinType,
		condition: condition,
		args:      args,
	})
	return m
}

// Where adds a WHERE clause with validation
func (m *Model) Where(column string, operator string, value interface{}) *Model {
	if !validOperators[strings.ToUpper(operator)] {
		panic(ErrInvalidOperator)
	}

	m.query.wheres = append(m.query.wheres, whereClause{
		column:   sanitizeColumn(column),
		operator: strings.ToUpper(operator),
		value:    value,
	})
	return m
}

// OrWhere adds an OR WHERE clause with validation
func (m *Model) OrWhere(column string, operator string, value interface{}) *Model {
	if !validOperators[strings.ToUpper(operator)] {
		panic(ErrInvalidOperator)
	}

	m.query.orWheres = append(m.query.orWheres, whereClause{
		column:   sanitizeColumn(column),
		operator: strings.ToUpper(operator),
		value:    value,
	})
	return m
}

// Get executes the query and returns all matching records
func (m *Model) Get() ([]map[string]interface{}, error) {
	query, args := m.buildSelectQuery()

	if m.db.queryLog {
		defer logQuery(query, args, time.Now())
	}

	stmt, err := m.prepareQuery(query)
	if err != nil {
		return nil, fmt.Errorf("prepare query error: %w", err)
	}

	rows, err := stmt.QueryContext(m.ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	return m.scanRows(rows)
}

// First returns the first matching record
func (m *Model) First() (map[string]interface{}, error) {
	m.query.limit = 1
	results, err := m.Get()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrNoRows
	}
	return results[0], nil
}

// Create inserts a new record with better error handling
func (m *Model) Create(data map[string]interface{}) (map[string]interface{}, error) {
	if len(data) == 0 {
		return nil, ErrInvalidValue
	}

	// Create a new map to avoid modifying the input map
	newData := make(map[string]interface{}, len(data)+2)
	for k, v := range data {
		newData[k] = v
	}

	// Add timestamps if they don't exist
	now := time.Now()
	if _, exists := newData["created_at"]; !exists {
		newData["created_at"] = now
	}
	if _, exists := newData["updated_at"]; !exists {
		newData["updated_at"] = now
	}

	columns := make([]string, 0, len(newData))
	values := make([]interface{}, 0, len(newData))
	placeholders := make([]string, 0, len(newData))

	i := 1
	for column, value := range newData {
		columns = append(columns, sanitizeColumn(column))
		values = append(values, value)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		i++
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING *",
		m.query.table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	if m.db.queryLog {
		defer logQuery(query, values, time.Now())
	}

	stmt, err := m.prepareQuery(query)
	if err != nil {
		return nil, fmt.Errorf("prepare query error: %w", err)
	}

	rows, err := stmt.QueryContext(m.ctx, values...)
	if err != nil {
		return nil, fmt.Errorf("create error: %w", err)
	}
	defer rows.Close()

	results, err := m.scanRows(rows)
	if err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no data returned after insert")
	}

	return results[0], nil
}

// Update updates matching records with improved error handling
func (m *Model) Update(data map[string]interface{}) (int64, error) {
	if len(data) == 0 {
		return 0, ErrInvalidValue
	}

	sets := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	i := 1
	for column, value := range data {
		sets = append(sets, fmt.Sprintf("%s = $%d", sanitizeColumn(column), i))
		values = append(values, value)
		i++
	}

	whereClause, whereValues := m.buildWhereClause(i)
	values = append(values, whereValues...)

	query := fmt.Sprintf(
		"UPDATE %s SET %s%s",
		m.query.table,
		strings.Join(sets, ", "),
		whereClause,
	)

	if m.db.queryLog {
		defer logQuery(query, values, time.Now())
	}

	stmt, err := m.prepareQuery(query)
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}

	result, err := stmt.ExecContext(m.ctx, values...)
	if err != nil {
		return 0, fmt.Errorf("update error: %w", err)
	}

	return result.RowsAffected()
}

// Delete deletes matching records with improved error handling
func (m *Model) Delete() (int64, error) {
	whereClause, values := m.buildWhereClause(1)

	query := fmt.Sprintf(
		"DELETE FROM %s%s",
		m.query.table,
		whereClause,
	)

	if m.db.queryLog {
		defer logQuery(query, values, time.Now())
	}

	stmt, err := m.prepareQuery(query)
	if err != nil {
		return 0, fmt.Errorf("prepare query error: %w", err)
	}

	result, err := stmt.ExecContext(m.ctx, values...)
	if err != nil {
		return 0, fmt.Errorf("delete error: %w", err)
	}

	return result.RowsAffected()
}

// Helper methods
func (m *Model) buildSelectQuery() (string, []interface{}) {
	var queryBuilder strings.Builder
	var values []interface{}

	queryBuilder.WriteString(fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(m.query.selections, ", "),
		m.query.table,
	))

	// Add joins
	for _, join := range m.query.joins {
		if join.condition != "" {
			queryBuilder.WriteString(fmt.Sprintf(" %s %s ON %s", join.joinType, join.table, join.condition))
			values = append(values, join.args...)
		} else {
			queryBuilder.WriteString(fmt.Sprintf(" %s %s", join.joinType, join.table))
		}
	}

	// Add where clauses
	whereClause, whereValues := m.buildWhereClause(len(values) + 1)
	queryBuilder.WriteString(whereClause)
	values = append(values, whereValues...)

	// Add group by
	if len(m.query.groupBy) > 0 {
		queryBuilder.WriteString(" GROUP BY ")
		queryBuilder.WriteString(strings.Join(m.query.groupBy, ", "))
	}

	// Add having
	if len(m.query.having) > 0 {
		queryBuilder.WriteString(" HAVING ")
		for i, having := range m.query.having {
			if i > 0 {
				queryBuilder.WriteString(" AND ")
			}
			queryBuilder.WriteString(having.condition)
			values = append(values, having.args...)
		}
	}

	// Add order by
	if m.query.orderBy != "" {
		queryBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s", m.query.orderBy, m.query.orderDir))
	}

	// Add limit and offset
	if m.query.limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d", m.query.limit))
	}
	if m.query.offset > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET %d", m.query.offset))
	}

	return queryBuilder.String(), values
}

func (m *Model) buildWhereClause(startIndex int) (string, []interface{}) {
	if len(m.query.wheres) == 0 && len(m.query.orWheres) == 0 {
		return "", nil
	}

	var whereBuilder strings.Builder
	var values []interface{}
	whereBuilder.WriteString(" WHERE ")

	paramIndex := startIndex

	for i, where := range m.query.wheres {
		if i > 0 {
			whereBuilder.WriteString(" AND ")
		}
		whereBuilder.WriteString(fmt.Sprintf("%s %s $%d", where.column, where.operator, paramIndex))
		values = append(values, where.value)
		paramIndex++
	}

	for i, orWhere := range m.query.orWheres {
		if len(m.query.wheres) > 0 || i > 0 {
			whereBuilder.WriteString(" OR ")
		}
		whereBuilder.WriteString(fmt.Sprintf("%s %s $%d", orWhere.column, orWhere.operator, paramIndex))
		values = append(values, orWhere.value)
		paramIndex++
	}

	return whereBuilder.String(), values
}

func (m *Model) scanRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}

		results = append(results, row)
	}

	return results, nil
}

func (m *Model) prepareQuery(query string) (*sql.Stmt, error) {
	m.db.mu.RLock()
	stmt, ok := m.db.prepared[query]
	m.db.mu.RUnlock()

	if ok {
		return stmt, nil
	}

	m.db.mu.Lock()
	defer m.db.mu.Unlock()

	// Double-check after acquiring write lock
	if stmt, ok = m.db.prepared[query]; ok {
		return stmt, nil
	}

	stmt, err := m.db.PrepareContext(m.ctx, query)
	if err != nil {
		return nil, err
	}

	m.db.prepared[query] = stmt
	return stmt, nil
}

func sanitizeColumn(column string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return -1
	}, column)
}

func sanitizeTableName(table string) string {
	return sanitizeColumn(table)
}

func sanitizeColumns(columns []string) []string {
	sanitized := make([]string, len(columns))
	for i, col := range columns {
		sanitized[i] = sanitizeColumn(col)
	}
	return sanitized
}

func logQuery(query string, args []interface{}, start time.Time) {
	duration := time.Since(start)
	fmt.Printf("[ORM] Query (%v):\n%s\nArgs: %v\n", duration, query, args)
}

// Cleanup closes all prepared statements
func (db *Orm) Cleanup() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	var errs []string
	for query, stmt := range db.prepared {
		if err := stmt.Close(); err != nil {
			errs = append(errs, fmt.Sprintf("failed to close statement for query %q: %v", query, err))
		}
	}

	db.prepared = make(map[string]*sql.Stmt)

	if len(errs) > 0 {
		return fmt.Errorf("cleanup errors: %s", strings.Join(errs, "; "))
	}
	return nil
}
