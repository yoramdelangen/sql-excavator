package mysql

import (
	sql "github.com/yoramdelangen/sql-excavator"
)

func mysqlOptions() *sql.SqlGrammar {
	return &sql.SqlGrammar{
		QuoteRune: '`',
    BindingPlaceholder: '?',


		// Clauses
		SelectClause:   []byte("select"),
		FromClause:     []byte("from"),
		WhereClause:    []byte("where"),
		AndClause:      []byte("and"),
		OrClause:       []byte("or"),
		PaginateClause: "limit %[2]d, %[1]d", // 1 = page, 2 = limit
	}
}

func Init() {
	sql.RegisterGrammar("mysql", mysqlOptions())
}
