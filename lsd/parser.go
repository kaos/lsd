//line parser.y:3
package main

import __yyfmt__ "fmt"

//line parser.y:4
import (
	"strings"
)

type OpFunc func(lhs, rhs interface{}) bool

//line parser.y:14
type lsdSymType struct {
	yys     int
	str     string
	strs    []string
	punct   rune
	num     float64
	hdrs    RequestHeaders
	arg     HeaderArgument
	opvalue interface{}
	opfunc  OpFunc
}

const GET = 57346
const LTEQ = 57347
const GTEQ = 57348
const EQI = 57349
const IDENT = 57350
const HEADER = 57351
const NUMBER = 57352
const PUNCT = 57353

var lsdToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"GET",
	"'<'",
	"'='",
	"'>'",
	"'!'",
	"LTEQ",
	"GTEQ",
	"EQI",
	"IDENT",
	"HEADER",
	"NUMBER",
	"PUNCT",
	"'~'",
}
var lsdStatenames = [...]string{}

const lsdEofCode = 1
const lsdErrCode = 2
const lsdMaxDepth = 200

//line parser.y:181
func wrapFilterFunc(op OpFunc, rhs interface{}) FilterFunc {
	return func(lhs interface{}) bool {
		return op(lhs, rhs)
	}
}

//line yacctab:1
var lsdExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const lsdNprod = 25
const lsdPrivate = 57344

var lsdTokenNames []string
var lsdStates []string

const lsdLast = 35

var lsdAct = [...]int{

	25, 13, 14, 15, 12, 17, 18, 19, 24, 13,
	14, 15, 16, 17, 18, 19, 27, 8, 11, 9,
	16, 21, 5, 26, 3, 2, 1, 22, 28, 10,
	20, 23, 7, 6, 4,
}
var lsdPact = [...]int{

	21, -1000, 12, -1000, 9, -1000, 5, -1000, -4, -1000,
	7, -1000, 4, -1000, -1000, -1000, -8, -1000, -1000, -1000,
	-1000, -1000, 11, -1000, -1000, 1, -1000, 11, 1,
}
var lsdPgo = [...]int{

	0, 34, 33, 32, 30, 18, 29, 27, 0, 26,
}
var lsdR1 = [...]int{

	0, 9, 1, 1, 2, 2, 3, 3, 3, 6,
	6, 5, 5, 5, 5, 5, 5, 5, 5, 4,
	4, 7, 7, 8, 8,
}
var lsdR2 = [...]int{

	0, 3, 0, 3, 0, 2, 1, 1, 3, 1,
	2, 1, 1, 1, 2, 1, 1, 1, 1, 1,
	1, 0, 2, 3, 1,
}
var lsdChk = [...]int{

	-1000, -9, 4, 12, -1, 13, -2, -3, 12, 14,
	-6, -5, 8, 5, 6, 7, 16, 9, 10, 11,
	-4, 14, -7, -5, 16, -8, 12, 15, -8,
}
var lsdDef = [...]int{

	0, -2, 0, 2, 1, 4, 3, 5, 6, 7,
	21, 9, 0, 11, 12, 13, 15, 16, 17, 18,
	8, 19, 20, 10, 14, 22, 24, 0, 23,
}
var lsdTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 8, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	5, 6, 7, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 16,
}
var lsdTok2 = [...]int{

	2, 3, 4, 9, 10, 11, 12, 13, 14, 15,
}
var lsdTok3 = [...]int{
	0,
}

var lsdErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	lsdDebug        = 0
	lsdErrorVerbose = false
)

type lsdLexer interface {
	Lex(lval *lsdSymType) int
	Error(s string)
}

type lsdParser interface {
	Parse(lsdLexer) int
	Lookahead() int
}

type lsdParserImpl struct {
	lookahead func() int
}

func (p *lsdParserImpl) Lookahead() int {
	return p.lookahead()
}

func lsdNewParser() lsdParser {
	p := &lsdParserImpl{
		lookahead: func() int { return -1 },
	}
	return p
}

const lsdFlag = -1000

func lsdTokname(c int) string {
	if c >= 1 && c-1 < len(lsdToknames) {
		if lsdToknames[c-1] != "" {
			return lsdToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func lsdStatname(s int) string {
	if s >= 0 && s < len(lsdStatenames) {
		if lsdStatenames[s] != "" {
			return lsdStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func lsdErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !lsdErrorVerbose {
		return "syntax error"
	}

	for _, e := range lsdErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + lsdTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := lsdPact[state]
	for tok := TOKSTART; tok-1 < len(lsdToknames); tok++ {
		if n := base + tok; n >= 0 && n < lsdLast && lsdChk[lsdAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if lsdDef[state] == -2 {
		i := 0
		for lsdExca[i] != -1 || lsdExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; lsdExca[i] >= 0; i += 2 {
			tok := lsdExca[i]
			if tok < TOKSTART || lsdExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if lsdExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += lsdTokname(tok)
	}
	return res
}

func lsdlex1(lex lsdLexer, lval *lsdSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = lsdTok1[0]
		goto out
	}
	if char < len(lsdTok1) {
		token = lsdTok1[char]
		goto out
	}
	if char >= lsdPrivate {
		if char < lsdPrivate+len(lsdTok2) {
			token = lsdTok2[char-lsdPrivate]
			goto out
		}
	}
	for i := 0; i < len(lsdTok3); i += 2 {
		token = lsdTok3[i+0]
		if token == char {
			token = lsdTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = lsdTok2[1] /* unknown char */
	}
	if lsdDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", lsdTokname(token), uint(char))
	}
	return char, token
}

func lsdParse(lsdlex lsdLexer) int {
	return lsdNewParser().Parse(lsdlex)
}

func (lsdrcvr *lsdParserImpl) Parse(lsdlex lsdLexer) int {
	var lsdn int
	var lsdlval lsdSymType
	var lsdVAL lsdSymType
	var lsdDollar []lsdSymType
	_ = lsdDollar // silence set and not used
	lsdS := make([]lsdSymType, lsdMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	lsdstate := 0
	lsdchar := -1
	lsdtoken := -1 // lsdchar translated into internal numbering
	lsdrcvr.lookahead = func() int { return lsdchar }
	defer func() {
		// Make sure we report no lookahead when not parsing.
		lsdstate = -1
		lsdchar = -1
		lsdtoken = -1
	}()
	lsdp := -1
	goto lsdstack

ret0:
	return 0

ret1:
	return 1

lsdstack:
	/* put a state and value onto the stack */
	if lsdDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", lsdTokname(lsdtoken), lsdStatname(lsdstate))
	}

	lsdp++
	if lsdp >= len(lsdS) {
		nyys := make([]lsdSymType, len(lsdS)*2)
		copy(nyys, lsdS)
		lsdS = nyys
	}
	lsdS[lsdp] = lsdVAL
	lsdS[lsdp].yys = lsdstate

lsdnewstate:
	lsdn = lsdPact[lsdstate]
	if lsdn <= lsdFlag {
		goto lsddefault /* simple state */
	}
	if lsdchar < 0 {
		lsdchar, lsdtoken = lsdlex1(lsdlex, &lsdlval)
	}
	lsdn += lsdtoken
	if lsdn < 0 || lsdn >= lsdLast {
		goto lsddefault
	}
	lsdn = lsdAct[lsdn]
	if lsdChk[lsdn] == lsdtoken { /* valid shift */
		lsdchar = -1
		lsdtoken = -1
		lsdVAL = lsdlval
		lsdstate = lsdn
		if Errflag > 0 {
			Errflag--
		}
		goto lsdstack
	}

lsddefault:
	/* default state action */
	lsdn = lsdDef[lsdstate]
	if lsdn == -2 {
		if lsdchar < 0 {
			lsdchar, lsdtoken = lsdlex1(lsdlex, &lsdlval)
		}

		/* look through exception table */
		xi := 0
		for {
			if lsdExca[xi+0] == -1 && lsdExca[xi+1] == lsdstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			lsdn = lsdExca[xi+0]
			if lsdn < 0 || lsdn == lsdtoken {
				break
			}
		}
		lsdn = lsdExca[xi+1]
		if lsdn < 0 {
			goto ret0
		}
	}
	if lsdn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			lsdlex.Error(lsdErrorMessage(lsdstate, lsdtoken))
			Nerrs++
			if lsdDebug >= 1 {
				__yyfmt__.Printf("%s", lsdStatname(lsdstate))
				__yyfmt__.Printf(" saw %s\n", lsdTokname(lsdtoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for lsdp >= 0 {
				lsdn = lsdPact[lsdS[lsdp].yys] + lsdErrCode
				if lsdn >= 0 && lsdn < lsdLast {
					lsdstate = lsdAct[lsdn] /* simulate a shift of "error" */
					if lsdChk[lsdstate] == lsdErrCode {
						goto lsdstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if lsdDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", lsdS[lsdp].yys)
				}
				lsdp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if lsdDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", lsdTokname(lsdtoken))
			}
			if lsdtoken == lsdEofCode {
				goto ret1
			}
			lsdchar = -1
			lsdtoken = -1
			goto lsdnewstate /* try again in the same state */
		}
	}

	/* reduction by production lsdn */
	if lsdDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", lsdn, lsdStatname(lsdstate))
	}

	lsdnt := lsdn
	lsdpt := lsdp
	_ = lsdpt // guard against "declared and not used"

	lsdp -= lsdR2[lsdn]
	// lsdp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if lsdp+1 >= len(lsdS) {
		nyys := make([]lsdSymType, len(lsdS)*2)
		copy(nyys, lsdS)
		lsdS = nyys
	}
	lsdVAL = lsdS[lsdp+1]

	/* consult goto table to find next state */
	lsdn = lsdR1[lsdn]
	lsdg := lsdPgo[lsdn]
	lsdj := lsdg + lsdS[lsdp].yys + 1

	if lsdj >= lsdLast {
		lsdstate = lsdAct[lsdg]
	} else {
		lsdstate = lsdAct[lsdj]
		if lsdChk[lsdstate] != -lsdn {
			lsdstate = lsdAct[lsdg]
		}
	}
	// dummy call; replaced with literal code
	switch lsdnt {

	case 1:
		lsdDollar = lsdS[lsdpt-3 : lsdpt+1]
		//line parser.y:40
		{
			req := lsdlex.(*lsdLex).Request
			req.table = lsdDollar[2].str
			req.headers = lsdDollar[3].hdrs
		}
	case 2:
		lsdDollar = lsdS[lsdpt-0 : lsdpt+1]
		//line parser.y:47
		{
			lsdVAL.hdrs = make(RequestHeaders, 0)
		}
	case 3:
		lsdDollar = lsdS[lsdpt-3 : lsdpt+1]
		//line parser.y:51
		{
			lsdVAL.hdrs = append(lsdDollar[1].hdrs,
				RequestHeader{
					name: lsdDollar[2].str,
					arg:  lsdDollar[3].arg,
				})
		}
	case 4:
		lsdDollar = lsdS[lsdpt-0 : lsdpt+1]
		//line parser.y:60
		{
			lsdVAL.arg = nil
		}
	case 5:
		lsdDollar = lsdS[lsdpt-2 : lsdpt+1]
		//line parser.y:64
		{
			switch arg := lsdDollar[1].arg.(type) {
			case nil:
				lsdVAL.arg = lsdDollar[2].arg
			case []HeaderArgument:
				lsdVAL.arg = append(arg, lsdDollar[2].arg)
			default:
				lsdVAL.arg = []HeaderArgument{lsdDollar[1].arg, lsdDollar[2].arg}
			}
		}
	case 6:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:76
		{
			lsdVAL.arg = HeaderArgument(lsdDollar[1].str)
		}
	case 7:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:80
		{
			lsdVAL.arg = HeaderArgument(lsdDollar[1].num)
		}
	case 8:
		lsdDollar = lsdS[lsdpt-3 : lsdpt+1]
		//line parser.y:84
		{
			lsdVAL.arg = HeaderArgument(
				&FilterExpression{
					field: lsdDollar[1].str,
					op:    wrapFilterFunc(lsdDollar[2].opfunc, lsdDollar[3].opvalue),
				})
		}
	case 9:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:93
		{
			lsdVAL.opfunc = lsdDollar[1].opfunc
		}
	case 10:
		lsdDollar = lsdS[lsdpt-2 : lsdpt+1]
		//line parser.y:97
		{
			tst := lsdDollar[2].opfunc
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return !tst(lhs, rhs)
			}
		}
	case 11:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:105
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return compare(lhs, rhs) == -1
			}
		}
	case 12:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:111
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return compare(lhs, rhs) == 0
			}
		}
	case 13:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:117
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return compare(lhs, rhs) == 1
			}
		}
	case 14:
		lsdDollar = lsdS[lsdpt-2 : lsdpt+1]
		//line parser.y:123
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return regexp_matchf(lhs, rhs, "i")
			}
		}
	case 15:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:129
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return regexp_match(lhs, rhs)
			}
		}
	case 16:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:135
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return compare(lhs, rhs) <= 0
			}
		}
	case 17:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:141
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return compare(lhs, rhs) >= 0
			}
		}
	case 18:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:147
		{
			lsdVAL.opfunc = func(lhs, rhs interface{}) bool {
				return comparef(lhs, rhs, CF_IGNORE_CASE) == 0
			}
		}
	case 19:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:154
		{
			lsdVAL.opvalue = lsdDollar[1].num
		}
	case 20:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:158
		{
			lsdVAL.opvalue = strings.Join(lsdDollar[1].strs, " ")
		}
	case 21:
		lsdDollar = lsdS[lsdpt-0 : lsdpt+1]
		//line parser.y:163
		{
			lsdVAL.strs = []string{}
		}
	case 22:
		lsdDollar = lsdS[lsdpt-2 : lsdpt+1]
		//line parser.y:167
		{
			lsdVAL.strs = append(lsdDollar[1].strs, lsdDollar[2].str)
		}
	case 23:
		lsdDollar = lsdS[lsdpt-3 : lsdpt+1]
		//line parser.y:172
		{
			lsdVAL.str = lsdDollar[1].str + string(lsdDollar[2].punct) + lsdDollar[3].str
		}
	case 24:
		lsdDollar = lsdS[lsdpt-1 : lsdpt+1]
		//line parser.y:176
		{
			lsdVAL.str = lsdDollar[1].str
		}
	}
	goto lsdstack /* stack new state and value */
}
