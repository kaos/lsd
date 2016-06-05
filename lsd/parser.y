// -*- mode: go -*-
%{

package main

import (
	"strings"
)

type OpFunc func (lhs, rhs interface{}) bool

%}

%union {
        str     string
	strs    []string
	punct   rune
        num     float64
        hdrs    RequestHeaders
        arg     HeaderArgument
	opvalue interface{}
	opfunc  OpFunc
}

%type <hdrs>    Headers
%type <arg>     HeaderArgument HeaderValue
%type <opvalue> Rhs
%type <opfunc>  Op OpExpr
%type <strs>    StringRest
%type <str>     NextString

%token GET '<' '=' '>' '!' LTEQ GTEQ EQI
%token <str> IDENT HEADER
%token <num> NUMBER
%token <punct> PUNCT

%%

Request: GET IDENT Headers
{
        req := lsdlex.(*lsdLex).Request
        req.table = $2
        req.headers = $3
}

Headers:
{
        $$ = make(RequestHeaders, 0)
}
| Headers HEADER HeaderArgument
{
        $$ = append($1,
                RequestHeader{
                        name: $2,
                        arg: $3,
                })
}

HeaderArgument:
{
        $$ = nil
}
| HeaderArgument HeaderValue
{
        switch arg := $1.(type) {
        case nil:
                $$ = $2
        case []HeaderArgument:
                $$ = append(arg, $2)
        default:
                $$ = []HeaderArgument{$1, $2}
        }
}

HeaderValue: IDENT
{
        $$ = HeaderArgument($1)
}
| NUMBER
{
        $$ = HeaderArgument($1)
}
| IDENT OpExpr Rhs
{
	$$ = HeaderArgument(
		&FilterExpression{
			field: $1,
			op: wrapFilterFunc($2, $3),
		})
}

OpExpr: Op
{
	$$ = $1
}
| '!' Op
{
	tst := $2
	$$ = func (lhs, rhs interface{}) bool {
		return !tst(lhs, rhs)
	}
}

Op: '<'
{
	$$ = func (lhs, rhs interface{}) bool {
		return compare(lhs, rhs) == -1
	}
}
| '='
{
	$$ = func (lhs, rhs interface{}) bool {
		return compare(lhs, rhs) == 0
	}
}
| '>'
{
	$$ = func (lhs, rhs interface{}) bool {
		return compare(lhs, rhs) == 1
	}
}
| '~' '~'
{
	$$ = func (lhs, rhs interface{}) bool {
		return regexp_matchf(lhs, rhs, "i")
	}
}
| '~'
{
	$$ = func (lhs, rhs interface{}) bool {
		return regexp_match(lhs, rhs)
	}
}
| LTEQ
{
	$$ = func (lhs, rhs interface{}) bool {
		return compare(lhs, rhs) <= 0
	}
}
| GTEQ
{
	$$ = func (lhs, rhs interface{}) bool {
		return compare(lhs, rhs) >= 0
	}
}
| EQI
{
	$$ = func (lhs, rhs interface{}) bool {
		return comparef(lhs, rhs, CF_IGNORE_CASE) == 0
	}
}

Rhs: NUMBER
{
	$$ = $1
}
| StringRest
{
	$$ = strings.Join($1, " ")
}

StringRest:
{
	$$ = []string{}
}
| StringRest NextString
{
	$$ = append($1, $2)
}

NextString: NextString PUNCT NextString
{
	$$ = $1 + string($2) + $3
}
| IDENT
{
	$$ = $1
}


%%

func wrapFilterFunc (op OpFunc, rhs interface{}) FilterFunc {
	return func (lhs interface{}) bool {
		return op(lhs, rhs)
	}
}
