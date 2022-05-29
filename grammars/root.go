package grammars

import (
	"fmt"
	"strings"
	str "strings"
)

// Inspiration:
// - https://github.com/laravel/framework/blob/d787947998089e6eac4883d90418601fdb5c529d/src/Illuminate/Database/Grammar.php
// - https://github.com/arthurkushman/buildsqlx/blob/master/builder.go -- Inspiration for a query builder setup.
// - https://github.com/doug-martin/goqu/tree/master/dialect -- How GO to solve this problem.
// - https://github.com/laravel/framework/tree/9.x/src/Illuminate/Database -- Query language for Mysql, SQL Server, PostgreSQL

var (
	StarClause string = "*"
	Nested     string = "(%s)"
	Space      rune   = ' '
	Separator  rune   = ','
)

type SqlGrammar struct {
	// TODO: setup setting
	// Rune used to escape table and columns.
	QuoteRune          rune
	BindingPlaceholder rune
	Operators          []rune

	// Simple clause names
	Select string
	From   string
	Where  string
	And    string
	Or     string

	Insert string
	Values string
	Update string
	Set    string

	Null      string
	IsNull    string
	IsNotNull string

	// PaginationClause; "LIMIT {limit} OFFSET {offset}"
	Paginate string

	// Enabled configuration
	EnableTableWrap  bool
	EnableColumnWrap bool
}

func (s *SqlGrammar) Table(b *str.Builder, name string) {
	b.WriteString(s.wrapTable(name))
}

func (s *SqlGrammar) AddColumn(b *str.Builder, column string) {
	if column == StarClause {
		b.WriteString(StarClause)
		return
	}

	// TODO: Check for functions usage.
	// TODO: Check for aliases/table prefixes
	b.WriteString(s.wrapColumn(column))
}

func (s *SqlGrammar) CompilePaginate(b *str.Builder, offset int, limit uint) {
	b.WriteRune(Space)
	b.WriteString(fmt.Sprintf(s.Paginate, offset, limit))
}

func (s *SqlGrammar) CompileWhere(b *str.Builder, where WhereClause) bool {
	hasBindings := false

	b.WriteString(s.wrapColumn(where.Column))
	b.WriteRune(Space)

	if where.Type == BindingType.NULL {
		b.WriteString(s.IsNull)
	} else if where.Type == BindingType.NOT_NULL {
		b.WriteString(s.IsNotNull)
	} else if where.Type == BindingType.RAW {
		// TODO: query raw
	} else {
		hasBindings = true
		b.WriteString(where.Operator)
		b.WriteRune(Space)

		// TODO: Numeric bindings
		b.WriteRune(s.BindingPlaceholder)
	}

	return hasBindings
}

// TODO: allowence for structs
func (s *SqlGrammar) CompileInsert(b *str.Builder, table string, data map[string]interface{}) (*str.Builder, []interface{}) {
	// get columns and values in seperated var's
	keys := make([]string, 0, len(data))
	values := make([]interface{}, 0, len(data))

	for k, v := range data {
		keys = append(keys, s.wrapColumn(k)) // TODO: wrap column name
		values = append(values, v)
	}

	// INSERT INTO {table}
	b.WriteString(s.Insert)
	b.WriteRune(Space)
	s.Table(b, table)
	b.WriteRune(Space)

	// (col1, col2, etc)
	cols := strings.Join(keys, string(Separator)+string(Space))
	b.WriteString(fmt.Sprintf(Nested, cols))

	// values
	b.WriteRune(Space)
	b.WriteString(s.Values)
	b.WriteRune(Space)

	// set bindings
	vals := strings.Join(GeneratePlaceholders(len(data), string(s.BindingPlaceholder)), string(Separator)+string(Space))
	b.WriteString(fmt.Sprintf(Nested, vals))

	return b, values
}

func (s *SqlGrammar) wrapColumn(column string) string {
	if s.EnableColumnWrap == false {
		return column
	}
	return s.wrapQuote(column)
}

func (s *SqlGrammar) wrapTable(table string) string {
	if s.EnableTableWrap == false {
		return table
	}
	return s.wrapQuote(table)
}

func (s *SqlGrammar) wrapQuote(val string) string {
	return string(s.QuoteRune) + val + string(s.QuoteRune)
}

func GeneratePlaceholders(size int, value string) []string {
	arr := make([]string, size)
	for i := 0; i < size; i++ {
		arr[i] = value
	}
	return arr
}
