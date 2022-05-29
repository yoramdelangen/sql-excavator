package grammars

func NewMysqlQuery() *SqlGrammar {
	return &SqlGrammar{
		QuoteRune:          '`',
		BindingPlaceholder: '?',

		// Clauses
		Select: "select",
		From:   "from",
		Where:  "where",
		And:    "and",
		Or:     "or",

		Insert: "insert into",
		Values: "values",
		Update: "update",
		Set:    "set",

		Null:      "null",
		IsNull:    "is null",
		IsNotNull: "is not null",

		Paginate: "limit %[2]d, %[1]d", // 1 = page, 2 = limit

		EnableTableWrap:  true,
		EnableColumnWrap: true,
	}
}
