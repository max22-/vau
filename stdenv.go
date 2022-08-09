package main

import (
	"errors"
	//"log"
)


func stdenv() environment {
	return map[string]Value {
		"abc": Number(2),
			"vau": proc(vau),
			"def": proc(def),
			"+": proc(add),
			"car": proc(car),
			"list": proc(list),
	}
}

func vau(args Value, e environment) (Value, error) {
	vars := args.(*Cell).car.(*Cell)
	call_env_sym := args.(*Cell).cdr.(*Cell).car.(Symbol)
	body := args.(*Cell).cdr.(*Cell).cdr.(*Cell).car
	return proc(func (args Value, call_env environment) (Value, error) {
		new_env := make(environment)
		for k, v := range e {
			new_env[k] = v
		}
		pvars := vars
		pargs := args
		for pvars != Nil {
			new_env[string(pvars.car.(Symbol))] = pargs.(*Cell).car
			pvars = pvars.cdr.(*Cell)
			pargs = pargs.(*Cell).cdr
		}
		new_env[string(call_env_sym)] = call_env
		return eval(body, new_env)
	}), nil
	
}

func def(args Value, e environment) (Value, error) {
	name := args.(*Cell).car.(Symbol)
	val, err := eval(args.(*Cell).cdr.(*Cell).car, e)
	if (err != nil) {
		return val, err
	}
	e[string(name)] = val
	return val, nil
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
