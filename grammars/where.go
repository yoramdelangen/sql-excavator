package grammars

type BooleanTypeName string
type BindingTypeName string

var BooleanType = struct {
	OR     BooleanTypeName
	AND    BooleanTypeName
	NESTED BooleanTypeName
}{
	OR:     "OR",
	AND:    "AND",
	NESTED: "NESTED",
}

var BindingType = struct {
	BASIC    BindingTypeName
	NULL     BindingTypeName
	NOT_NULL BindingTypeName
	RAW      BindingTypeName
}{
	BASIC:    "BASIC",
	NULL:     "NULL",
	NOT_NULL: "NOT_NULL",
	RAW:      "RAW",
}

type WhereClause struct {
	Type        BindingTypeName
	BooleanType BooleanTypeName
	Column      string
	Operator    []byte
}

type Where struct {
	bindings []interface{}

	// When type is nested we can have multiple wheres as nested
	wheres []WhereClause
}

func NewWhere() Where {
  return Where{}
}

func (w Where) GetClauses() []WhereClause {
	return w.wheres
}

func (w Where) GetBindings() []interface{} {
	return w.bindings
}

func (w Where) And(column string, operator []byte, val interface{}) Where {
	w.wheres = append(w.wheres, WhereClause{
		Type:        BindingType.BASIC,
		BooleanType: BooleanType.AND,
		Column:      column,
		Operator:    operator,
	})
	w.bindings = append(w.bindings, val)
	return w
}

func (w Where) Or(column string, operator []byte, val interface{}) Where {
	w.wheres = append(w.wheres, WhereClause{
		Type:        BindingType.BASIC,
		BooleanType: BooleanType.OR,
		Column:      column,
		Operator:    operator,
	})
	w.bindings = append(w.bindings, val)
	return w
}
func (w Where) Null(column string, boolean BooleanTypeName) Where {
	w.wheres = append(w.wheres, WhereClause{
		Type:        BindingType.NULL,
		BooleanType: boolean,
		Column:      column,
		Operator:    []byte(""),
	})

	return w
}

func (w Where) NotNull(column string, boolean BooleanTypeName) Where {
	w.wheres = append(w.wheres, WhereClause{
		Type:        BindingType.NOT_NULL,
		BooleanType: boolean,
		Column:      column,
		Operator:    []byte(""),
	})

	return w
}

func (w Where) Raw(raw string, boolean BooleanTypeName) Where {
	w.wheres = append(w.wheres, WhereClause{
		Type:        BindingType.RAW,
		BooleanType: boolean,
		Column:      "",
		Operator:    []byte(raw),
	})

	return w
}
