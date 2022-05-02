package sql_excavator

import (
	"bytes"
	"fmt"
	str "strings"

	"github.com/yoramdelangen/sql-excavator/grammars"
)

func NewBuilder(dialect string) *Builder {
	return &Builder{
		grammar: *GetGrammar(dialect),
	}
}

type Builder struct {
	grammar  SqlGrammar
	table    string
	wheres   []grammars.Where
	join     []grammars.Join
	columns  []string
	orderBy  []string
	bindings []interface{}
	limit    uint
	offset   int
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func (r *Builder) reset() {
	r.table = ""
	r.wheres = []grammars.Where{}
	r.join = []grammars.Join{}
	r.columns = []string{StarClause}
	r.orderBy = []string{}
	r.offset = -1
}

// Sql function will render a SQL string based upon the dilect.
// It will follow a certain structure for building up the query.
// The content of this function will probably be moved to the
// grammar file because the order of operations can change per dialece.
func (r *Builder) Sql() (string, []interface{}) {
	var b *str.Builder = &str.Builder{}

	b.Write(r.grammar.SelectClause)
	// TODO: Adding columns
	r.createSelect(b)
	r.createFrom(b)

	// set TABLE
	r.grammar.Table(b, r.table)

	// setup the where columns
	r.createWhere(b)

	// TODO: Set WHERE clauses
	if r.limit != 0 && r.offset > -1 {
		r.grammar.Paginate(b, r.offset, r.limit)
	}

	return b.String(), r.bindings
}

// Basically the start for a new query. On every call it will call function `reset` to reset the current instance.
// In this way we can safely buildup and have a clean builder instance everytim/
// Heavily inspired by the following libraries: Laravel Eloquent, gopu
func (r *Builder) Table(table string) *Builder {
	r.reset()
	r.table = table

	return r
}

// Allows to select certain columns within the query. Each argument is string and will be set.
// Calling this function will set the columns and flushes previously set columns.
// When conditionally adding columns use `AddSelect` function.
func (r *Builder) Select(columns ...string) *Builder {
	r.columns = columns

	return r
}

// Conditionally add columns to the instance. Differents with `Select` this only appends.
// It will not flush the columns slice.
func (r *Builder) AddSelect(columns ...string) *Builder {
	r.columns = append(r.columns, columns...)
	return r
}

// Add a simple WHERE clause statement. WHen having only 2 arguments (include column) it will equally
// be to equal sign and generates with operator equal. Otherwise it will use the operator
// given as second argument and value as third argument.
// When having more then 3 total arguments it will panic.
func (r *Builder) Where(column string, args ...interface{}) *Builder {
	if len(args) > 2 {
		panic("Where cannot have more then 3 arguments")
	}

	value := args[0]
	operator := []byte("=")

	if len(args) == 2 {
		operator = []byte(args[0].(string))
		value = args[1]
	}

	r.wheres = append(r.wheres, grammars.Where{}.Or(column, operator, value))

	return r
}

func (r *Builder) OrWhere(column string, args ...interface{}) *Builder {
	fmt.Printf("Where Or Type: %#v %T\n", args, args[2])

	value := args[0]
	operator := []byte("=")

	if len(args) == 2 {
		operator = []byte(args[0].(string))
		value = args[1]
	}

	r.wheres = append(r.wheres, grammars.Where{}.Or(column, operator, value))

	return r
}

// Helper function to generate pagination like feature. In the majority of dialects it will generate a LIMIT and OFFSET.
func (r *Builder) Paginate(page uint, limit uint) *Builder {
	r.limit = limit
	r.offset = int((page - 1) * limit)

	return r
}

// Create the FROM clause
func (r *Builder) createFrom(b *str.Builder) {
	b.WriteRune(Space)
	b.Write(r.grammar.FromClause)
	b.WriteRune(Space)
}

// Internal function for adding columns to to grammar state and write it into the strings Builder.
func (r *Builder) createSelect(b *str.Builder) {
	if len(r.columns) == 0 {
		r.columns = []string{StarClause}
	}

	columnLength := len(r.columns)
	for i, column := range r.columns {
		b.WriteRune(Space)

		r.grammar.AddColumn(b, column)

		// are there more columns?
		if i < (columnLength - 1) {
			b.WriteRune(Separator)
		}
	}
}

func (r *Builder) createWhere(b *str.Builder) {
	if len(r.wheres) == 0 {
		return
	}

	b.WriteRune(Space)
	b.Write(r.grammar.WhereClause)

	x := 0
	for _, where := range r.wheres {
		fmt.Printf("WHERE: %+v\n", where)
		// Lets start with the basic type and only add with AND
		for i, clause := range where.GetClauses() {
			if i == 0 && x > 0 {
				b.WriteRune(Space)

				if bytes.Equal(clause.WhereType, grammars.WhereType.AND) {
					b.Write(r.grammar.AndClause)
				} else if bytes.Equal(clause.WhereType, grammars.WhereType.OR) {
					b.Write(r.grammar.OrClause)
				}
			}

			b.WriteRune(Space)
			// TODO: should wrap the column
			b.WriteString(clause.Column)
			b.WriteRune(Space)
			b.Write(clause.Operator)
			b.WriteRune(Space)

			// Write placeholder..
			b.WriteRune(r.grammar.BindingPlaceholder)

			r.bindings = append(r.bindings, where.GetBindings()...)
		}
		x++
	}
	fmt.Println("")
	fmt.Println("")

	fmt.Printf("QueryBuilder: %#v\n\n", r.bindings)
}
