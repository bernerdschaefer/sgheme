package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		sicp = flag.Bool("sicp", false, "enter directly to sicp evaluator")
	)

	flag.Parse()

	env := &environment{
		values: map[object]object{
			scmSymbol("true"):  TRUE,
			scmSymbol("false"): FALSE,
		},
	}

	for _, proc := range makePrimitives(env) {
		env.values[scmSymbol(proc.name)] = proc
	}

	currentScanner = newScanner(os.Stdin)

	if *sicp {
		env.values[scmSymbol("load-file!")].(primitive).Call(&cell{
			car: scmString("sicp.lisp"),
			cdr: NIL,
		})

		return
	}

	for {
		fmt.Printf("> ")
		fmt.Printf("%s\n", eval(read(), env))
	}
}
