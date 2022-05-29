package mysql

import (
	sql "github.com/yoramdelangen/sql-excavator"
)

func mysqlOptions() *sql.SqlGrammar {
	return &sql.SqlGrammar{
		QuoteRune: '`',
    BindingPlaceholder: '?',

		// Clauses
		SelectClause:   "select",
		FromClause:     "from",
		WhereClause:    "where",
		AndClause:      "and",
		OrClause:       "or",
		PaginateClause: "limit %[2]d, %[1]d", // 1 = page, 2 = limit

    NullValue: "null",
    IsNullValue: "is null",
    IsNotNullValue: "is not null",
	}
}

func Init() {
	sql.RegisterGrammar("mysql", mysqlOptions())
}
