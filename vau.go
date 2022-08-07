package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"github.com/chzyer/readline"
)

type Sexp interface {
	isSexp()
	definition
}

type Atom interface {
	isAtom()
	Sexp
}

type Symbol string
type Number int
type List []Sexp

func (_ Symbol) isSexp() {}
func (_ Symbol) isAtom() {}
func (s Symbol) String() string { return string(s) }

func (_ Number) isSexp() {}
func (_ Number) isAtom() {}
func (n Number) String() string { return strconv.Itoa(int(n)) }

func (_ List) isSexp() {}
func (l List) String() string {
	es := make([]string, len(l))
	for i, e := range l {
		es[i] = fmt.Sprintf("%v", e)
	}
	return "(" + strings.Join(es, " ") + ")"
}

func tokenize(line string) []string {
	line = strings.ReplaceAll(line, "(", " ( ")
	line = strings.ReplaceAll(line, ")", " ) ")
	return strings.Fields(line)
}

func parseAtom(tokens []string, idx int) (Atom, int, error) {
	token := tokens[idx]
	if n, err := strconv.Atoi(token); err == nil {
		return Number(n), idx + 1, nil
	} else {
		return Symbol(token), idx + 1, nil
	}	
}

func parseList(tokens []string, idx int) (l List, ridx int, err error) {
	var sexp Sexp
	ridx = idx + 1 // We match the opening parenthese
	for ridx < len(tokens) {
		if tokens[ridx] == ")" {
			ridx += 1
			return
		}
		sexp, ridx, err = parse(tokens, ridx)
		if err != nil {
			return nil, idx, err
		}
		l = append(l, sexp)
	}
	return nil, idx, errors.New("Unterminated s-expression")
}

func parse(tokens []string, idx int) (sexp Sexp, ridx int, err error) {
	if len(tokens) == 0 {
		return nil, idx, nil
	}
	if idx >= len(tokens) {
		return nil, idx, errors.New("EOF")
	}
	token := tokens[idx]
	if token == ")" {
		return nil, idx, errors.New("Syntax error (unexpected '(')")
	} else if token == "(" {
		return parseList(tokens, idx)
	} else {
		return parseAtom(tokens, idx)
	}
}

type definition interface {
	isDefinition()
	String() string
}

type proc func(e environment, args Sexp) (Sexp, error)
func proc_to_func(p proc) func(environment, Sexp) (Sexp, error) {
	return func(environment, Sexp) (Sexp, error) (p)
}

func (p proc) String() string {
	return "<proc " + fmt.Sprintf("%v", proc_to_func(p)) + ">"
}

func (_ Symbol) isDefinition() {}
func (_ Number) isDefinition() {}
func (_ List) isDefinition() {}
func (_ proc) isDefinition() {}

type environment map[string]definition

func eval(x Sexp, e environment) (definition, error) {
	switch x.(type) {
	case Symbol:
		if def, ok := e[string(x.(Symbol))]; ok {
			return def, nil
		} else {
			return x, nil
		}
	case List:
		x0 := x.(List)[0]
		p, err := eval(x0, e)
		if err != nil {
			return p, err
		}
		switch p.(type) {
		case proc:
			return proc_to_func(p.(proc))(e, x.(List)[1:])
		default:
			return nil, errors.New(fmt.Sprintf("%v", x0) + " is not a procedure")
		}
	default:
		return x, nil
	}
}

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	program := "(begin (define r 10) (* pi r r))"
	sexp, _, err := parse(tokenize(program), 0)
	if err == nil {
		fmt.Println(sexp)
	} else {
		fmt.Printf("error: %v\n", err)
	}
	//return

	env := stdenv()

	// REPL
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		sexp, _, err := parse(tokenize(line), 0)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}
		val, err := eval(sexp, env)
		if err == nil {
			fmt.Println(val)
		} else {
			fmt.Printf("error: %v\n", err);
		}
	}
}
