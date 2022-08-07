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
}

type Atom interface {
	isAtom()
	Sexp
	String() string
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

func main() {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	program := "(begin (define r 10) (* pi r r))"
	sexp, _, err := parse(tokenize(program), 0)
	if err == nil {
		fmt.Printf("%q\n", sexp)
	} else {
		fmt.Printf("error: %v\n", err)
	}
	//return


	// REPL
	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		sexp, idx, err := parse(tokenize(line), 0)
		if err == nil {
			fmt.Printf("idx = %d\n", idx)
			fmt.Printf("%q\n", sexp)
		} else {
			fmt.Printf("error: %v\n", err);
		}
	}
}
