package sql_excavator

import (
	"strings"

	"github.com/yoramdelangen/sql-excavator/grammars"
)


func GetGrammar(name string) *grammars.SqlGrammar {
	dialect := strings.ToLower(name)

	switch dialect {
	case "mysql":
		return grammars.NewMysqlQuery()
	}

  panic("Dialect not found")
}
