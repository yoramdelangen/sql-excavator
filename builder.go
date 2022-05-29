package sql_excavator

import (
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
	grammar  grammars.SqlGrammar
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
	r.columns = []string{grammars.StarClause}
	r.orderBy = []string{}
	r.offset = -1
}

// Sql function will render a SQL string based upon the dilect.
// It will follow a certain structure for building up the query.
// The content of this function will probably be moved to the
// grammar file because the order of operations can change per dialece.
func (r *Builder) Sql() (string, []interface{}) {
	var b *str.Builder = &str.Builder{}

	b.WriteString(r.grammar.Select)
	// TODO: Adding columns
	r.createSelect(b)
	r.createFrom(b)

	// set TABLE
	r.grammar.Table(b, r.table)

	// setup the where columns
	r.createWhere(b)

	// TODO: Set WHERE clauses
	if r.limit != 0 && r.offset > -1 {
		r.grammar.CompilePaginate(b, r.offset, r.limit)
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
	operator := "="

	if len(args) == 2 {
		operator = args[0].(string)
		value = args[1]
	}

	r.wheres = append(r.wheres, grammars.Where{}.And(column, operator, value))

	return r
}

func (r *Builder) OrWhere(column string, args ...interface{}) *Builder {
	value := args[0]
	operator := "="

	if len(args) == 2 {
		operator = args[0].(string)
		value = args[1]
	}

	r.wheres = append(r.wheres, grammars.Where{}.Or(column, operator, value))

	return r
}

// Set when column is null
func (r *Builder) WhereNull(column string) *Builder {
	r.wheres = append(r.wheres, grammars.Where{}.Null(column, grammars.BooleanType.AND))
	return r
}

// Or when column is null
func (r *Builder) OrWhereNull(column string) *Builder {
	r.wheres = append(r.wheres, grammars.Where{}.Null(column, grammars.BooleanType.OR))
	return r
}

// Set when column is null
func (r *Builder) WhereNotNull(column string) *Builder {
	r.wheres = append(r.wheres, grammars.Where{}.NotNull(column, grammars.BooleanType.AND))
	return r
}

// Or when column is null
func (r *Builder) OrWhereNotNull(column string) *Builder {
	r.wheres = append(r.wheres, grammars.Where{}.NotNull(column, grammars.BooleanType.OR))
	return r
}

// Helper function to generate pagination like feature. In the majority of dialects it will generate a LIMIT and OFFSET.
func (r *Builder) Paginate(page uint, limit uint) *Builder {
	r.limit = limit
	r.offset = int((page - 1) * limit)

	return r
}

func (r *Builder) Insert(args map[string]interface{}) (string, []interface{}) {
	var b *str.Builder = &str.Builder{}

  b, bindings := r.grammar.CompileInsert(b, r.table, args)

	return b.String(), bindings
}

// Create the FROM clause
func (r *Builder) createFrom(b *str.Builder) {
	b.WriteRune(grammars.Space)
	b.WriteString(r.grammar.From)
	b.WriteRune(grammars.Space)
}

// Internal function for adding columns to to grammar state and write it into the strings Builder.
func (r *Builder) createSelect(b *str.Builder) {
	if len(r.columns) == 0 {
		r.columns = []string{grammars.StarClause}
	}

	columnLength := len(r.columns)
	for i, column := range r.columns {
		b.WriteRune(grammars.Space)

		r.grammar.AddColumn(b, column)

		// are there more columns?
		if i < (columnLength - 1) {
			b.WriteRune(grammars.Separator)
		}
	}
}

func (r *Builder) createWhere(b *str.Builder) {
	if len(r.wheres) == 0 {
		return
	}

	b.WriteRune(grammars.Space)
	b.WriteString(r.grammar.Where)

	x := 0
	for _, where := range r.wheres {
		// Lets start with the basic type and only add with AND
		for i, clause := range where.GetClauses() {
			if i == 0 && x > 0 {
				b.WriteRune(grammars.Space)

				if clause.BooleanType == grammars.BooleanType.AND {
					b.WriteString(r.grammar.And)
				} else if clause.BooleanType == grammars.BooleanType.OR {
					b.WriteString(r.grammar.Or)
				}
			}

			b.WriteRune(grammars.Space)

			hasBindings := r.grammar.CompileWhere(b, clause)

			// continue when type was RAW
			if hasBindings == false {
				continue
			}

			r.bindings = append(r.bindings, where.GetBindings()...)
		}
		x++
	}
}
