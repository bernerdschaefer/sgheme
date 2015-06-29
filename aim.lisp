; <bernerd>
(define (map proc args)
  (if (null? args)
    '()
    (cons (proc (car args)) (map proc (cdr args)))))
(define apply-in-underlying-scheme apply)
(define (newline) (display ""))
(define (cond?) false)

; shorthand
; from http://www.lispworks.com/documentation/HyperSpec/Body/f_car_c.htm
(define (caar x)        (car (car x)))
(define (cadr x)        (car (cdr x)))
(define (cdar x)        (cdr (car x)))
(define (cddr x)        (cdr (cdr x)))
(define (caaar x)       (car (car (car x))))
(define (caadr x)       (car (car (cdr x))))
(define (cadar x)       (car (cdr (car x))))
(define (caddr x)       (car (cdr (cdr x))))
(define (cdaar x)       (cdr (car (car x))))
(define (cdadr x)       (cdr (car (cdr x))))
(define (cddar x)       (cdr (cdr (car x))))
(define (cdddr x)       (cdr (cdr (cdr x))))
(define (caaaar x)      (car (car (car (car x)))))
(define (caaadr x)      (car (car (car (cdr x)))))
(define (caadar x)      (car (car (cdr (car x)))))
(define (caaddr x)      (car (car (cdr (cdr x)))))
(define (cadaar x)      (car (cdr (car (car x)))))
(define (cadadr x)      (car (cdr (car (cdr x)))))
(define (caddar x)      (car (cdr (cdr (car x)))))
(define (cadddr x)      (car (cdr (cdr (cdr x)))))
(define (cdaaar x)      (cdr (car (car (car x)))))
(define (cdaadr x)      (cdr (car (car (cdr x)))))
(define (cdadar x)      (cdr (car (cdr (car x)))))
(define (cdaddr x)      (cdr (car (cdr (cdr x)))))
(define (cddaar x)      (cdr (cdr (car (car x)))))
(define (cddadr x)      (cdr (cdr (car (cdr x)))))
(define (cdddar x)      (cdr (cdr (cdr (car x)))))
(define (cddddr x)      (cdr (cdr (cdr (cdr x)))))

; end bernerd


(define (eval exp env)
  (cond ((self-evaluating? exp) exp)
        ((variable? exp) (lookup-variable-value exp env))
        ((quoted? exp) (text-of-quotation exp))
        ((assignment? exp) (eval-assignment exp env))
        ((definition? exp) (eval-definition exp env))
        ((if? exp) (eval-if exp env))
        ((lambda? exp)
         (make-procedure
           (lambda-parameters exp)
           (lambda-body exp)
           env))
        ((begin? exp)
         (eval-sequence
           (begin-actions exp)
           env))
        ((cond? exp)
         (eval (cond->if exp) env))
        ((application? exp)
         (apply (eval (operator exp) env)
                (list-of-values
                  (operands exp)
                  env)))
        (else
          (error "Unknown expression type: EVAL" exp))))

(define (apply procedure arguments)
  (cond ((primitive-procedure? procedure)
         (apply-primitive-procedure
           procedure
           arguments))
        ((compound-procedure? procedure)
         (eval-sequence
           (procedure-body procedure)
           (extend-environment
             (procedure-parameters procedure)
             arguments
             (procedure-environment procedure))))
        (else
          (error "Unknown procedure type: APPLY" procedure))))

(define (list-of-values exps env)
  (if (no-operands? exps)
    '()
    (cons (eval (first-operand exps) env)
          (list-of-values (rest-operands exps) env))))

(define (eval-if exp env)
  (if (true? (eval (if-predicate exp) env))
    (eval (if-consequant exp) env)
    (eval (if-alternative exp) env)))

; eval a sequence of expressions in the environment
; and return the value of the last expression.
(define (eval-sequence exps env)
  (cond ((last-exp? exps)
         (eval (first-exp exps) env))
        (else
          (eval (first-exp exps) env)
          (eval-sequence (rest-exps exps) env))))

(define (eval-assignment exp env)
  (set-variable-value!
    (assignment-variable exp)
    (eval (assignment-value exp) env)
    env)
  'ok)

(define (eval-definition exp env)
  (define-variable!
    (definition-variable exp)
    (eval (definition-value exp) env)
    env)
  'ok)

;;;;;;;;;;;;;;;;

; syntax definition

(define (self-evaluating? exp)
  (cond ((number? exp) exp)
        ((string? exp) exp)
        (else false)))

(define (variable? exp) (symbol? exp))

; quotation: (quote <text-of-quotation>)
(define (quoted? exp)
  (tagged-list? exp 'quote))

(define (text-of-quotation exp)
  (cadr exp))

(define (tagged-list? exp tag)
  (if (pair? exp)
    (eq? (car exp) tag)
    false))

; (set! <var> <value>)
(define (assignment? exp)
  (tagged-list? exp 'set!))

(define (assignment-variable exp)
  (cadr exp))

(define (assignment-value exp)
  (caddr exp))

; definitions
; (define <var> <value>)
; (define (<var> <param1> ... <paramN>) <body>)
; the latter is syntactic sugar for:
;   (define <var> (lambda (<param1> ... <paramN>) <body>))

(define (definition? exp)
  (tagged-list? exp 'define))

(define (definition-variable exp)
  (if (symbol? (cadr exp))
    (cadr exp)
    (caadr exp)))

(define (definition-value exp)
  (if (symbol? (cadr exp))
    (caddr exp)
    (make-lambda
      (cdadr exp)    ; parameters
      (cddr exp))))  ; body

; lambda
(define (lambda? exp)
  (tagged-list? exp 'lambda))

(define (lambda-parameters exp)
  (cadr exp))

(define (lambda-body exp)
  (cddr exp))

(define (make-lambda parameters body)
  (cons 'lambda (cons parameters body)))

; conditionals

; (if <predicate> <consequant> <alternative>)
; (if <predicate> <consequant>) ; => 'false
(define (if? exp) (tagged-list? exp 'if))
(define (if-predicate exp) (cadr exp))
(define (if-consequant exp) (caddr exp))
(define (if-alternative exp)
  (if (not (null? (cdddr)))
    (cadddr exp)
    'false))

(define (make-if predicate
                 consequant
                 alternative)
  (list 'if
        predicate
        consequant
        alternative))

; evaluate expressions left to right and return the last value
; (begin (display "1") '2) ; => 2
(define (begin? exp)
  (tagged-list? exp 'begin))
(define (begin-actions exp)
  (cdr exp))
(define (last-exp? seq) (null? (cdr seq)))
(define (first-exp seq) (car seq))
(define (rest-exps seq) (cdr seq))

; convert sequence into a single expression
; (sequence->exp seq)
(define (sequence->exp seq)
  (cond ((null? seq) seq)
        ((last-exp? seq) (first-exp seq))
        (else (make-begin seq))))

; procedure application
(define (application? exp) (pair? exp))
(define (operator exp) (car exp))
(define (operands exp) (cdr exp))
(define (no-operands? ops) (null? ops))
(define (first-operand ops) (car ops))
(define (rest-operands ops) (cdr ops))

; cond ...

;;;;;;;;;;;;;;;;;;;;;
;; data structures ;;
;;;;;;;;;;;;;;;;;;;;;

; anything not false is true
(define (true? x)
  (not (eq? x false)))

(define (false? x)
  (eq? x false))

; representing procedures

; assumes:
; (apply-primitive-procedure ⟨proc⟩ ⟨args⟩)
; (primitive-procedure? ⟨proc⟩)
(define (make-procedure parameters body env)
  (list 'procedure parameters body env))
(define (compound-procedure? exp)
  (tagged-list? exp 'procedure))
(define (procedure-parameters exp)
  (cadr exp))
(define (procedure-body exp)
  (caddr exp))
(define (procedure-environment exp)
  (cadddr exp))

; environments

(define (enclosing-environment env) (cdr env))
(define (first-frame env) (car env))
(define the-empty-environment '())

; each frame is represented as a pair of lists.
(define (make-frame variables values)
  (cons variables values))
(define (frame-variables frame) (car frame))
(define (frame-values frame) (cdr frame))
(define (add-binding-to-frame! var val frame)
  (set-car! frame (cons var (frame-variables frame)))
  (set-cdr! frame (cons val (frame-values frame))))

; TODO: error checking
(define (extend-environment vars vals base-env)
  (cons (make-frame vars vals) base-env))

;
(define (lookup-variable-value var env)
  (if (eq? env the-empty-environment)
    (error "no variable found" var)
    (lookup-frame-variable-value
      var
      (frame-variables (first-frame env))
      (frame-values (first-frame env))
      (enclosing-environment env))))

(define (lookup-frame-variable-value var vars vals next-env)
  (cond ((null? vars) (lookup-variable-value var next-env))
        ((eq? var (car vars)) (car vals))
        (else (lookup-frame-variable-value var (cdr vars) (cdr vals) next-env))))

; scan for variable and change when found
(define (set-variable-value! var val env))

; scan first frame for existing variable
; and change value as in set-variable-value!,
; otherwise use add-binding-to-frame!
(define (define-variable! var val env)
  (add-binding-to-frame! var val (first-frame env)))

;;;;;;;;;;;;;;;;;;;;;;;
;; run the evaluator ;;
;;;;;;;;;;;;;;;;;;;;;;;

(define (setup-environment)
  ((lambda (initial-env)
     (define-variable! 'true true initial-env)
     (define-variable! 'false false initial-env)
     initial-env)
   (extend-environment
     (primitive-procedure-names)
     (primitive-procedure-objects)
     the-empty-environment)))

; primitive operations
(define (primitive-procedure? proc)
  (tagged-list? proc 'primitive))
(define (primitive-implementation proc)
  (cadr proc))

(define primitive-procedures
  (list (list 'car car)
        (list 'cdr cdr)
        (list 'cons cons)
        (list 'null? null?)
        (list 'set-car! set-car!)
        (list 'set-cdr! set-cdr!)
        (list 'eq? eq?)
        (list 'symbol? symbol?)
        (list 'number? number?)
        (list 'string? string?)
        (list 'not not)
        (list 'read read)
        (list 'display display)))

(define (primitive-procedure-names)
  (map car primitive-procedures))
(define (primitive-procedure-objects)
  (map (lambda (proc)
         (list 'primitive (cadr proc)))
       primitive-procedures))

; application
; we'll need to (define apply-in-underlying-scheme apply)
; before evaluating the metacircular apply
(define (apply-primitive-procedure proc args)
  (apply-in-underlying-scheme
    (primitive-implementation proc) args))

(define input-prompt ";;; M-Eval input:")
(define output-prompt ";;; M-Eval output:")

(define (driver-loop)
  (prompt-for-input input-prompt)
  (announce
    output-prompt
    (eval (read) the-global-environment))
  (driver-loop))

(define (announce prompt output)
  (announce-output prompt)
  (user-print output))

(define (prompt-for-input string)
  (newline) (newline)
  (display string) (newline))

(define (announce-output string)
  (newline) (display string) (newline))

; compound procedures include their environment,
; so specialize printing for them.
(define (user-print object)
  (if (compound-procedure? object)
    (display (list 'compound-procedure
                   (procedure-parameters object)
                   (procedure-body object)
                   '<procedure-env>))
    (display object)))

; run it!
(define the-global-environment
  (setup-environment))

(driver-loop)
