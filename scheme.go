package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var DEBUG = false

func eval(exp object, env *environment) object {
	if DEBUG {
		fmt.Printf("evaluating: %s\n", exp)
	}

	for {
		switch e := exp.(type) {
		case scmNumber:
			return exp
		case scmString:
			return exp
		case scmSymbol:
			return lookupVariableValue(exp, env)
		case *procedure:
			return exp
		case *cell:
			exp = evalCell(e, env)
			if f, ok := exp.(*tailCall); ok {
				exp = f.f
				env = f.env
				continue
			}
			return exp
		default:
			return raiseError("Unknown expression type: EVAL", exp)
		}
	}
}

func evalCell(e *cell, env *environment) object {
	switch e.car {
	case scmSymbol("quote"):
		return car(e.cdr)
	case scmSymbol("set!"):
		set(car(cdr(e)), car(cdr(cdr(e))), env)
		return OK
	case scmSymbol("define"):
		define(e, env)
		return OK
	case scmSymbol("lambda"):
		return &procedure{
			parameters: car(cdr(e)),
			body:       cdr(cdr(e)),
			env:        env,
		}
	case scmSymbol("cond"):
		return evalCond(e, env)
	case scmSymbol("if"):
		var (
			predicate   = car(cdr(e))
			consequant  = car(cdr(cdr(e)))
			alternative = car(cdr(cdr(cdr(e))))
		)

		if eval(predicate, env) != FALSE {
			return eval(consequant, env)
		}

		if alternative == NIL {
			return FALSE
		}

		return eval(alternative, env)

	case scmSymbol("with-error-handler"):
		// (with-error-handler
		//   (lambda (err) (display err))
		//   (lambda () (eval (read))))
		return withErrorHandler(
			eval(car(cdr(e)), env),
			eval(car(cdr(cdr(e))), env),
		)

	default: // application
		return apply(
			eval(e.car, env),
			evalList(e.cdr, env),
		)
	}
}

func apply(operator, arguments object) object {
	if DEBUG {
		fmt.Printf("[prim] applying %v with: %s\n", operator, arguments)
	}

	switch op := operator.(type) {
	case primitive:
		return op.Call(arguments)

	case *procedure:
		return evalSequence(
			op.body,
			extendEnvironment(op.parameters, arguments, op.env),
		)

	default:
		return raiseError("Unknown expression type: APPLY", operator)
	}
}

func evalCond(e object, env *environment) object {
	clauses := cdr(e)

	for ; clauses != NIL; clauses = cdr(clauses) {
		var (
			predicate = car(car(clauses))
			actions   = cdr(car(clauses))
		)

		switch {
		case predicate == scmSymbol("else"):
			return evalSequence(actions, env)
		case eval(predicate, env) != FALSE:
			return evalSequence(actions, env)
		}
	}

	return FALSE
}

type tailCall struct {
	f   object
	env *environment
}

// forceTailCall evaluates o immediately if it is a tail call.
// Otherwise it returns the object unmodified.
func forceTailCall(o object) object {
	if f, ok := o.(*tailCall); ok {
		return eval(f.f, f.env)
	}

	return o
}

// evalSequence evaluates the list of objects in exps.
// The last item in the list is returned
// as a special tailCall value which the evaluator inlines.
func evalSequence(exps object, env *environment) object {
	for {
		if cdr(exps) == NIL {
			return &tailCall{car(exps), env}
		}

		eval(car(exps), env)
		exps = cdr(exps)
	}
}

func evalList(o object, env *environment) object {
	if o == NIL {
		return NIL
	}

	return cons(
		eval(car(o), env),
		evalList(cdr(o), env),
	)
}

// withErrorHandler calls thunk and returns its value
// unless there's an error, in which case
// it calls handler and returns its value.
//
// Note that tail-calls are forced to keep a stack frame,
// so this should not be called recursively.
func withErrorHandler(handler, thunk object) (retval object) {
	defer func() {
		if err := recover(); err != nil {
			if sErr, ok := err.(scmError); ok {
				retval = forceTailCall(apply(handler, &cell{sErr, NIL}))
				return
			}

			panic(err)
		}
	}()

	return forceTailCall(apply(thunk, NIL))
}

type object interface{}

type cell struct {
	car object
	cdr object
}

func (c *cell) String() string {
	parts := []string{}

	for ; c != nil && c.car != nil; c = c.cdr.(*cell) {
		parts = append(parts, fmt.Sprint(c.car))
	}

	return "(" + strings.Join(parts, " ") + ")"
}

type scmNumber float64

type scmString string

type scmSymbol string

// read waits to read a complete expression from r
// and returns it.
func read() object {
	e := currentScanner.read()
	if e == EOF {
		os.Exit(0)
	}
	return e
}

func set(variable, val object, env *environment) object {
	val = eval(val, env)

	for ; env != nil; env = env.outer {
		if _, ok := env.values[variable]; ok {
			env.values[variable] = val
			return OK
		}
	}

	return raiseError("variable not in environment", variable)
}

func define(definition object, env *environment) object {
	var (
		variable = car(cdr(definition))
		value    object
	)

	if symbolP(variable) == TRUE {
		value = car(cdr(cdr(definition)))
	} else {
		variable = car(variable)

		parameters := cdr(car(cdr(definition)))
		body := cdr(cdr(definition))

		value = cons(
			scmSymbol("lambda"),
			cons(parameters, body),
		)
	}

	env.values[variable] = eval(value, env)
	return OK
}

func car(o object) object {
	if o == NIL {
		return NIL
	}

	if cell, ok := o.(*cell); ok {
		return cell.car
	}

	return raiseError("Not a list", o)
}

func cdr(o object) object {
	if o == NIL {
		return NIL
	}

	if cell, ok := o.(*cell); ok {
		return cell.cdr
	}

	return raiseError("Not a list", o)
}

func cons(o, list object) object {
	return &cell{o, list}
}

var (
	NIL = &cell{}

	TRUE  scmSymbol = "true"
	FALSE scmSymbol = "false"

	OK = scmSymbol("ok")
)

func symbolP(o object) object {
	switch o.(type) {
	case scmSymbol:
		return TRUE
	default:
		return FALSE
	}
}

func nullP(o object) object {
	if o == NIL {
		return TRUE
	}

	return FALSE
}

func setCar(list, o object) object {
	if list == NIL {
		return raiseError("Cannot set car on empty list")
	}

	list.(*cell).car = o
	return OK
}

func setCdr(list, o object) object {
	if list == NIL {
		return raiseError("Cannot set cdr on empty list")
	}

	list.(*cell).cdr = o
	return OK
}

func eq(a, b object) object {
	if a == b {
		return TRUE
	}

	return FALSE
}

func not(o object) object {
	if o == FALSE {
		return TRUE
	}

	return FALSE
}

//////

type environment struct {
	values map[object]object
	outer  *environment
}

func lookupVariableValue(o object, env *environment) object {
	if value, ok := env.values[o]; ok {
		return value
	}

	if env.outer == nil {
		return raiseError("Undefined variable in environment", o)
	}

	return lookupVariableValue(o, env.outer)
}

func extendEnvironment(parameters, arguments object, env *environment) *environment {
	env = &environment{
		values: map[object]object{},
		outer:  env,
	}

	for ; parameters != NIL && arguments != NIL; parameters, arguments = cdr(parameters), cdr(arguments) {
		env.values[car(parameters)] = car(arguments)
	}

	return env
}

//////

type procedure struct {
	parameters object
	body       object
	env        *environment
}

func (p *procedure) String() string {
	return fmt.Sprintf("(lambda %s %s)", p.parameters, p.body)
}

//////

// scmError is currently opaque from within the scheme environment.
type scmError error

func raiseError(objects ...object) object {
	var s []string

	for _, o := range objects {
		s = append(s, fmt.Sprint(o))
	}

	panic(errors.New(strings.Join(s, " | ")))
}
