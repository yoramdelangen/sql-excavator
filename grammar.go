package sql_excavator

import (
	"fmt"
	str "strings"
)

// Inspiration:
// - https://github.com/laravel/framework/blob/d787947998089e6eac4883d90418601fdb5c529d/src/Illuminate/Database/Grammar.php
// - https://github.com/arthurkushman/buildsqlx/blob/master/builder.go -- Inspiration for a query builder setup.
// - https://github.com/doug-martin/goqu/tree/master/dialect -- How GO to solve this problem.
// - https://github.com/laravel/framework/tree/9.x/src/Illuminate/Database -- Query language for Mysql, SQL Server, PostgreSQL

var (
  StarClause string = "*"
  Space rune = ' '
  Separator rune = ','
)

type SqlGrammar struct {
  // TODO: setup setting
  // Rune used to escape table and columns.
  QuoteRune rune

  // Whats the "SELECT" named.
  SelectClause []byte

  // Whats the "FROM" named.
  FromClause []byte

  // PaginationClause; "LIMIT {limit} OFFSET {offset}"
  PaginateClause string
}

func (s *SqlGrammar) Table(b *str.Builder, name string) {
  s.wrapQuote(b, name)
}

func (s *SqlGrammar) AddColumn(b *str.Builder, column string){
  if column == StarClause {
    b.WriteString(StarClause)
    return
  }

  // TODO: Check for functions usage.
  // TODO: Check for aliases/table prefixes
  s.wrapQuote(b, column)
}

func (s *SqlGrammar) Where(column string, args ...interface{}) {
  // WHERE statement
  fmt.Printf("Where column %s has arguments %s", column, args)
}

func (s *SqlGrammar) Paginate(b *str.Builder, offset int, limit uint) {
  b.WriteRune(Space)
  b.WriteString(fmt.Sprintf(s.PaginateClause, offset, limit))
}

func (s *SqlGrammar) wrapQuote(b *str.Builder, val string) {
  b.WriteRune(s.QuoteRune)
  b.WriteString(val)
  b.WriteRune(s.QuoteRune)
}
