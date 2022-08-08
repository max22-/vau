package main

import (
	"errors"
	//"log"
)


func stdenv() environment {
	return map[string]Value {
		"abc": Number(2),
			"+": proc(add),
			"car": proc(car),
			"list": proc(list),
	}
}

func add(args Value, e environment) (Value, error) {
	res := 0
	if args == Nil {
		return Number(res), nil
	}
	for {
		val, err := eval(args.(*Cell).car, e)
		if err != nil {
			return val, err
		}
		switch val.(type) {
		case Number:
			res += int(val.(Number))
		default:
			return nil, errors.New("+: expected numbers")
		}
		if args.(*Cell).cdr != Nil {
			args = args.(*Cell).cdr.(*Cell)
		} else {
			break
		}
	}
	return Number(res), nil
}

func car(args Value, e environment) (Value, error) {
	if args.(*Cell).cdr != Nil {
		return nil, errors.New("car: invalid number of arguments")
	} else if args.(*Cell).car == Nil {
		return nil, errors.New("car: cannot get car of empty list")
	} else {
		val, err := eval(args.(*Cell).car, e)
		if err != nil {
			return val, err
		}
		switch val.(type)  {
		case *Cell:
			return val.(*Cell).car, nil
		default:
			return nil, errors.New("car: expected cell")
		}
	}
}

func list(args Value, e environment) (Value, error) {
	res := args
	for args != Nil {
		if val, err := eval(args.(*Cell).car, e); err == nil {
			args.(*Cell).car = val
			args = args.(*Cell).cdr.(*Cell)
		} else {
			return val, err
		}
	}
	return res, nil
}
