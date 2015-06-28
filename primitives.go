package main

import "fmt"

func prim_debug(e object) object {
	DEBUG = true
	return OK
}

func prim_cdr(e object) object {
	return cdr(car(e))
}

func prim_apply(e object) object {
	return apply(car(e), car(cdr(e)))
}

func prim_null(e object) object {
	return nullP(car(e))
}

func prim_symbol(e object) object {
	return symbolP(car(e))
}

func prim_number(e object) object {
	if _, ok := e.(scmNumber); ok {
		return TRUE
	}
	return FALSE
}

func prim_string(e object) object {
	if _, ok := e.(scmString); ok {
		return TRUE
	}
	return FALSE
}

func prim_pair(e object) object {
	if c, ok := e.(*cell); ok && c != NIL {
		return TRUE
	}
	return FALSE
}

func prim_display(e object) object {
	fmt.Printf("%s\n", car(e))
	return OK
}

func prim_eq(e object) object {
	var (
		a = car(e)
		b = car(cdr(e))
	)

	return eq(a, b)
}

func prim_not(e object) object {
	if e == FALSE {
		return TRUE
	}

	return FALSE
}

func prim_cons(e object) object {
	return cons(car(e), car(cdr(e)))
}

func prim_car(e object) object {
	return car(car(e))
}

func prim_set_car(e object) object {
	return setCar(car(e), car(cdr(e)))
}

func prim_set_cdr(e object) object {
	return setCdr(car(e), car(cdr(e)))
}

func prim_list(e object) object {
	return e
}

func prim_error(e object) object {
	return raiseError(car(e))
}

func prim_read(_ object) object {
	return read()
}

type primitive struct {
	name string
	f    func(object) object
}

func (p primitive) String() string {
	return p.name
}

func (p primitive) Call(o object) object {
	return p.f(o)
}

func makePrimitives(env *environment) []primitive {
	return []primitive{
		{"debug!", prim_debug},
		{"cdr", prim_cdr},
		{"apply", prim_apply},
		{"null?", prim_null},
		{"symbol?", prim_symbol},
		{"number?", prim_number},
		{"string?", prim_string},
		{"pair?", prim_pair},
		{"read", prim_read},
		{"display", prim_display},
		{"eq?", prim_eq},
		{"not", prim_not},
		{"cons", prim_cons},
		{"car", prim_car},
		{"set-car!", prim_set_car},
		{"set-cdr!", prim_set_cdr},
		{"list", prim_list},
		{"error", prim_error},
		{"eval", func(o object) object {
			return eval(car(o), env)
		}},
	}
}