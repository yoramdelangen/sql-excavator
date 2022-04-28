package mysql_test

import (
  "fmt"
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
  qs := builder.Table("testing").Sql()

  equals(t, qs, "TestSimpleSelect", "select * from `testing`")
}

func TestSelectSingleColumn(t *testing.T) {
  qs := builder.Table("testing").Select("id").Sql()

  equals(t, qs, "TestSelectSingleColumn", "select `id` from `testing`")
}

func TestSelectMultipleColumns(t *testing.T) {
  qs := builder.
    Table("testing").
    Select("id", "title").
    Sql()

  equals(t, qs, "TestSelectMultipleColumns", "select `id`, `title` from `testing`")
}

func TestAddSelectSingleColumn(t *testing.T) {
  qs := builder.Table("testing").AddSelect("id").Sql()

  equals(t, qs, "TestAddSelectSingleColumn", "select *, `id` from `testing`")
}

func TestPaginate(t *testing.T) {
  // Build query
  qs := builder.Table("testing").Paginate(1, 100).Sql()

  equals(t, qs, "TestPaginate", "select * from `testing` LIMIT 100, 0")
}

func TestSimpleWhere(t *testing.T) {
  qs := builder.Table("testing").
    Where("column", "!=", true).
    Where("column", "ABC").
    Sql()
  fmt.Println(qs)
}
