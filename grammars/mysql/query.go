package mysql

import (
  sql "github.com/yoramdelangen/sql-excavator"
)

func mysqlOptions() *sql.SqlGrammar {
	return &sql.SqlGrammar{
    QuoteRune: '`',

    // Clauses
    SelectClause: []byte("select"),
    FromClause: []byte("from"),
    PaginateClause: "LIMIT %[2]d, %[1]d", // 1 = page, 2 = limit
  }
}

func Init() {
	sql.RegisterGrammar("mysql", mysqlOptions())
}
