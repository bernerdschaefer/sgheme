# SGHEME

sgheme is a toy implementation of scheme in Go.

I started it to:

  * understand meta-circular evaluators
  * see what implementing scheme in Go might look like

With that in mind, then,
sgheme provides a host environment
sufficiently powerful to explore
the meta-circular evaluators of
[The Roots of Lisp] and [SICP],
and not much else.

It was mostly written over a weekend,
first by transcribing the SICP meta-circular evaluator,
and then using that as an spec for the go version,
like the original LISP, except not in machine code.

My primary focus is currently on
exploring the concepts behind scheme
and meta-circular evaluators.

If you read the go code, you'll know this to be true.

But despite the slapdash design,
it's been surprisingly easy
to add new features and fix outstanding bugs
in support of writing scheme.

## Running

```
go build
./sgheme
```

You'll be dropped into sgheme's REPL.

If you pass `-sicp`, you'll instead
be dropped directly into the SICP evaluator
defined in [sicp.lisp].
It's a shortcut for `(load-file! "sicp.lisp")`.

  [sicp.lisp]: sicp.lisp

## References

A partial list of things I've been reading through for inspiration.

  - [The Roots of Lisp]
  - [SICP]
  - [A Micro-Manual for LISP - Not the Whole Truth][micro-manual]
  - [The Art of the Interpreter][AIM-453]
  - [SRFI]
  - Many random things from the [ContentCreationWiki][c2].

  [The Roots of Lisp]: http://www.paulgraham.com/rootsoflisp.html
  [SICP]: http://sarabander.github.io/sicp/
  [micro-manual]: http://www.cse.sc.edu/~mgv/csce330f13/micromanualLISP.pdf
  [AIM-453]: http://repository.readscheme.org/ftp/papers/ai-lab-pubs/AIM-453.pdf
  [SRFI]: http://srfi.schemers.org/
  [c2]: http://c2.com/cgi/wiki

## TODO

  - [ ] read more of SICP and create more TODOs
  - [ ] implement `cond` in sicp lisp
  - [x] tail-call optimization
  - [x] not broken scanner
  - [x] error handling without killing the process

## License

sgheme is Copyright (c) 2015 Bernerd Schaefer.
It is free software, and may be redistributed
under the terms specified in the [LICENSE] file.

  [LICENSE]: LICENSE.md
