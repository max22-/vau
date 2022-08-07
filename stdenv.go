package main

import "errors"

func stdenv() environment {
	return map[string]definition {
		"abc": Number(2),
			"+": proc(add),
			"car": proc(car),
	}
}

func add(e environment, args Sexp) (Sexp, error) {
	res := 0
	for _, arg := range args.(List) {
		e_arg, err := eval(arg, e)
		if err != nil {
			return nil, err
		}
		switch e_arg.(type) {
		case Number:
			res += int(e_arg.(Number))
		default:
			return nil, errors.New("+: expected numbers")
		}
	}
	return Number(res), nil
}

func car(_ environment, args Sexp) (Sexp, error) {
	if len(args.(List)) != 1 {
		return nil, errors.New("car: invalid number of arguments")
	}
	arg := args.(List)[0]
	switch arg.(type) {
	case List:
		if len(arg.(List)) < 1 {
			return nil, errors.New("car: cannot get car of empty list")
		}
		return Sexp(arg.(List)[0]), nil
	default:
		return nil, errors.New("car: expected list")
	}
}
