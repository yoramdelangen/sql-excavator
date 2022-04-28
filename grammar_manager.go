package sql_excavator

import (
	"strings"
	"sync"
)

var (
	dialects   = make(map[string]*SqlGrammar)
	dialectsMu sync.RWMutex
)

// type SqlGrammarManager struct {
//   grammars map[string]SqlGrammar
// }

// func (g &SqlGrammarManager) Register(name string, grammar *SqlGrammar) {
//   g.grammars[name] = grammar
// }

func RegisterGrammar(name string, grammar *SqlGrammar) {
	dialectsMu.Lock()
	defer dialectsMu.Unlock()

	lowerName := strings.ToLower(name)
	dialects[lowerName] = grammar
}

func GetGrammar(name string) *SqlGrammar {
	lowerName := strings.ToLower(name)

	// TODO: check if dialect exists in the list of dialects

	return dialects[lowerName]
}
