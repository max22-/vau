package main

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"github.com/chzyer/readline"
	//"log"
)

// Value: Atom, Cell, or builtin function
type Value interface {
	isValue()
	String() string
}

type Cell struct {
	car Value
	cdr Value
	Value
}

type Atom interface {
	isAtom()
	Value
}

type Symbol string
type Number int

var Nil *Cell = nil

type environment map[string]Value
type proc func(args Value, e environment) (Value, error)

func (_ Symbol) isValue() {}
func (_ Symbol) isAtom() {}
func (s Symbol) String() string { return string(s) }

func (_ Number) isValue() {}
func (_ Number) isAtom() {}
func (n Number) String() string { return strconv.Itoa(int(n)) }

func (_ Cell) isValue() {}
func (c *Cell) String() (res string) {
	if c == Nil {
		return "()"
	}
	res = "("
	for {
		if c.car == Nil {
			res += "()"
		} else {
			res += fmt.Sprintf("%v", c.car)
		}
		if c.cdr == Nil {
			res += ")"
			return
		} else {
			switch c.cdr.(type) {
			case Atom:
				res += " . " + fmt.Sprintf("%v", c.cdr) + ")"
				return
			case *Cell:
				if c.cdr == Nil {
					res += ")"
					return
				} else {
					res += " "
					c = c.cdr.(*Cell)
				}
			}
		}
	}
}


func (p proc) String() string {
	return "<proc " + fmt.Sprintf("%v", from_proc(p)) + ">"
}

func (_ proc) isValue() {}

func from_proc(p proc) func(Value, environment) (Value, error) {
	return func(Value, environment) (Value, error) (p)
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

func parseList(tokens []string, idx int) (head *Cell, ridx int, err error) {
	var val Value
	var c *Cell = Nil
	ridx = idx + 1 // We match the opening bracket
	for ridx < len(tokens) {
		if tokens[ridx] == ")" {
			ridx += 1
			return
		}
		val, ridx, err = parse(tokens, ridx)
		if err != nil {
			return Nil, idx, err
		}
		// append
		if head == Nil {
			head = new(Cell)
			c = head
			c.car = val
			c.cdr = Nil
		} else {
			c.cdr = new(Cell)
			c.cdr.(*Cell).car = val
			c.cdr.(*Cell).cdr = Nil
			c = c.cdr.(*Cell)
		}
	}
	return Nil, idx, errors.New("Unterminated s-expression")
}

func parse(tokens []string, idx int) (val Value, ridx int, err error) {
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

func eval(x Value, e environment) (Value, error) {
	switch x.(type) {
	case Symbol:
		if def, ok := e[string(x.(Symbol))]; ok {
			return def, nil
		} else {
			return x, nil
		}
	case *Cell:
		if x.(*Cell) == Nil {
			return Nil, nil
		}
		x0 := x.(*Cell).car
		p, err := eval(x0, e)
		if err != nil {
			return p, err
		}
		switch p.(type) {
		case proc:
			return from_proc(p.(proc))(x.(*Cell).cdr.(*Cell), e)
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
