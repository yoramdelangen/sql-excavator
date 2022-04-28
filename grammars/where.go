package grammars

type WhereTypeName []byte

var WhereType = struct {
	OR     WhereTypeName
	AND    WhereTypeName
	NESTED WhereTypeName
}{
	OR:     []byte("OR"),
	AND:    []byte("AND"),
	NESTED: []byte("NESTED"),
}

type WhereClause struct {
	WhereType WhereTypeName
	Column    string
	Operator  []byte
}

type Where struct {
	bindings []interface{}

	// When type is nested we can have multiple wheres as nested
	wheres []WhereClause
}

func (w Where) GetClauses() []WhereClause {
	return w.wheres
}
func (w Where) GetBindings() []interface{} {
	return w.bindings
}

func (w Where) And(column_ string, operator_ []byte, val interface{}) Where {
	w.wheres = append(w.wheres, WhereClause{
		WhereType: WhereType.AND,
		Column:    column_,
		Operator:  operator_,
	})
	w.bindings = append(w.bindings, val)
	return w
}

func (w Where) Or(column_ string, operator_ []byte, val interface{}) Where {
	w.wheres = append(w.wheres, WhereClause{
		WhereType: WhereType.OR,
		Column:    column_,
		Operator:  operator_,
	})
	w.bindings = append(w.bindings, val)
	return w
}
