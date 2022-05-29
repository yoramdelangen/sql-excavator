package mysql_test

import (
	"testing"

	sql "github.com/yoramdelangen/sql-excavator"
	"github.com/yoramdelangen/sql-excavator/grammars/mysql"
)

// enabling colors for testing: 'brew install grc' with guide: https://stackoverflow.com/questions/27242652/colorizing-golang-test-run-output

var (
	builder *sql.Builder
)

func equals(t *testing.T, query string, label string, expecting string) {
	if query != expecting {
		t.Errorf("Expected query for '%s'\nEXPECT: %+s\nGOT: %+s", label, expecting, query)
	}

	t.Logf("Method: %s Query: %s\n", label, query)
}

func TestMain(t *testing.M) {
	mysql.Init()
	builder = sql.NewBuilder("mysql")
	t.Run()
}

func TestSimpleSelect(t *testing.T) {
	// Build query
	qs, _ := builder.Table("testing").Sql()

	equals(t, qs, "TestSimpleSelect", "select * from `testing`")
}

func TestSelectSingleColumn(t *testing.T) {
	qs, _ := builder.Table("testing").Select("id").Sql()

	equals(t, qs, "TestSelectSingleColumn", "select `id` from `testing`")
}

func TestSelectMultipleColumns(t *testing.T) {
	qs, _ := builder.
		Table("testing").
		Select("id", "title").
		Sql()

	equals(t, qs, "TestSelectMultipleColumns", "select `id`, `title` from `testing`")
}

func TestAddSelectSingleColumn(t *testing.T) {
	qs, _ := builder.Table("testing").AddSelect("id").Sql()

	equals(t, qs, "TestAddSelectSingleColumn", "select *, `id` from `testing`")
}

func TestPaginate(t *testing.T) {
	// Build query
	qs, _ := builder.Table("testing").Paginate(1, 100).Sql()

	equals(t, qs, "TestPaginate", "select * from `testing` limit 100, 0")
}

func TestSimpleWhere(t *testing.T) {
	qs, _ := builder.Table("testing").
		Where("column", "!=", true).
		OrWhere("column", "ABC").
		Sql()

	equals(t, qs, "TestSimpleWhere", "select * from `testing` where column != ? or column = ?")
}

func TestSimpleWhereNull(t *testing.T) {
	qs, _ := builder.Table("testing").
		WhereNull("column").
		OrWhereNull("column2").
		Sql()

	equals(t, qs, "TestSimpleWhere", "select * from `testing` where column is null or column2 is null")
}
func TestSimpleWhereNullWithNotNull(t *testing.T) {
	qs, _ := builder.Table("testing").
		WhereNotNull("column").
		WhereNull("column2").
		Sql()

	equals(t, qs, "TestSimpleWhere", "select * from `testing` where column is not null and column2 is null")
}
