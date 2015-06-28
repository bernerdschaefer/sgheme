package main

import (
	"fmt"
	"os"
)

func main() {
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

	for {
		fmt.Printf("> ")
		fmt.Printf("%s\n", eval(read(), env))
	}
}
