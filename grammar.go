package sql_excavator

import (
	"fmt"
	str "strings"

	"github.com/yoramdelangen/sql-excavator/grammars"
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
	SelectClause string
	FromClause   string
	WhereClause  string
	AndClause    string
	OrClause     string

	NullValue      string
	IsNullValue    string
	IsNotNullValue string

	// PaginationClause; "LIMIT {limit} OFFSET {offset}"
	PaginateClause string
}

func (s *SqlGrammar) Table(b *str.Builder, name string) {
	s.wrapQuote(b, name)
}

func (s *SqlGrammar) AddColumn(b *str.Builder, column string) {
	if column == StarClause {
		b.WriteString(StarClause)
		return
	}

	// TODO: Check for functions usage.
	// TODO: Check for aliases/table prefixes
	s.wrapQuote(b, column)
}

func (s *SqlGrammar) Paginate(b *str.Builder, offset int, limit uint) {
	b.WriteRune(Space)
	b.WriteString(fmt.Sprintf(s.PaginateClause, offset, limit))
}

func (s *SqlGrammar) CompileWhere(b *str.Builder, where grammars.WhereClause) bool {
  hasBindings := false

	// TODO: should wrap the column
	b.WriteString(where.Column)
	b.WriteRune(Space)

	if where.Type == grammars.BindingType.NULL {
		b.WriteString(s.IsNullValue)
	} else if where.Type == grammars.BindingType.NOT_NULL {
		b.WriteString(s.IsNotNullValue)
	} else if where.Type == grammars.BindingType.RAW {
    // TODO: query raw
	} else {
    hasBindings = true
    b.Write(where.Operator)
    b.WriteRune(Space)
    // TODO: Numeric bindings
    b.WriteRune(s.BindingPlaceholder)
  }

  return hasBindings
}

func (s *SqlGrammar) wrapQuote(b *str.Builder, val string) {
	b.WriteRune(s.QuoteRune)
	b.WriteString(val)
	b.WriteRune(s.QuoteRune)
}
