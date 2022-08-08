package main

import "errors"

func stdenv() environment {
	return map[string]Value {
		"abc": Number(2),
			"+": proc(add),
			"car": proc(car),
	}
}

func add(args *Cell, e environment) (Value, error) {
	res := 0
	if args == nil {
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
		if args.cdr != nil {
			args = args.cdr.(*Cell)
		} else {
			break
		}
	}
	return Number(res), nil
}

func car(args *Cell, _ environment) (Value, error) {
	if args == nil {
		return nil, errors.New("car: cannot get car of empty list")
	} else {
		return args.car.(*Cell).car, nil
	}
}
