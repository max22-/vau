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

func add(args *Cell, e environment) (Value, error) {
	res := 0
	if args == Nil {
		return Number(res), nil
	}
	for {
		val, err := eval(args.car, e)
		if err != nil {
			return val, err
		}
		switch val.(type) {
		case Number:
			res += int(val.(Number))
		default:
			return nil, errors.New("+: expected numbers")
		}
		if args.cdr != Nil {
			args = args.cdr.(*Cell)
		} else {
			break
		}
	}
	return Number(res), nil
}

func car(args *Cell, e environment) (Value, error) {
	if args.cdr != Nil {
		return nil, errors.New("car: invalid number of arguments")
	} else if args.car == Nil {
		return nil, errors.New("car: cannot get car of empty list")
	} else {
		val, err := eval(args.car, e)
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

func list(args *Cell, e environment) (Value, error) {
	res := args
	for args != Nil {
		if val, err := eval(args.car, e); err == nil {
			args.car = val
			args = args.cdr.(*Cell)
		} else {
			return val, err
		}
	}
	return res, nil
}
