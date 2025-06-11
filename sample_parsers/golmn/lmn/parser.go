package lmn

import (
	"math"
)

type LmnParser struct {
	idx     int
	buf     []byte
	anchors map[string]any
}

func NewLmn() LmnParser {
	return LmnParser{
		idx:     0,
		buf:     nil,
		anchors: nil,
	}
}

func (lp *LmnParser) value() (any, error) {
	lp.skip()

	var val any
	var err error = nil

	if val, err = lp.getAnchorValue(); err == nil { // 캡쳐 가지오기 성공
		// Capture cannot be captured
		return val, nil
	}

	err = nil

	switch here := lp.here(); here {
	case '(':
		val, err = lp.dictionary(true)
	case '[':
		val, err = lp.list()
	case '\'', '"':
		val, err = lp.string()
	case '?':
		lp.idx++
		val = nil
	case '!':
		lp.idx++
		val = math.NaN()
	case '+', '-':
		lp.idx++
		next := lp.here()

		var sign int

		if here == '+' {
			sign = 1
		} else {
			sign = -1
		}

		// 숫자일 때
		if '0' <= next && next <= '9' {
			val, err = lp.number(sign)
		} else if next == '^' { // 무한일 때
			lp.idx++
			val = math.Inf(sign)
		} else {
			val = here == '+'
		}
	case '^':
		lp.idx++
		val, err = math.Inf(1), nil
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		val, err = lp.number(1)
	default:
		val, err = nil, lp.err(unexpectedTokenErr)
	}

	// ~ 캡쳐 처리
	lp.skip()

	if lp.here() == '~' {
		lp.idx++
		lp.skip()

		if cap, err := lp.anchor(); err != nil {
			return nil, err
		} else {
			if _, exist := lp.anchors[cap]; exist {
				return nil, lp.err(duplicatedCapErr)
			} else {
				lp.anchors[cap] = val
			}
		}
	}

	return val, err
}

func (lp *LmnParser) Parse(lmn string) (any, error) {
	lp.idx = 0
	lp.buf = []byte(lmn)
	lp.anchors = map[string]any{}

	var res any
	var err error

	// try parse top-level dictionary
	res, err = lp.topLevelDictionary()

	if err == nil { // top-level dictionary 성공
		return res, nil
	}

	// top-level dictionary 실패했을 때
	lp.idx = 0 // 인덱스 초기화

	res, err = lp.value()

	if err != nil {
		return nil, err
	}

	lp.skip()

	if lp.notEnd() {
		if lp.here() == ',' { // top-level list
			var topList = []any{res}
			lp.idx++
			lp.skip()

			for lp.notEnd() {
				if val, err := lp.value(); err != nil {
					return nil, err
				} else {
					topList = append(topList, val)
					lp.skip()
				}

				if here := lp.here(); here == ',' {
					lp.idx++
					lp.skip()
				} else if lp.end() {
					break
				} else {
					return nil, lp.err(unexpectedTokenErr)
				}
			}
			res = topList
		} else {
			lp.err(unexpectedTokenErr)
		}
	}

	return res, nil
}

func Parse(lmn string) (any, error) {
	p := NewLmn()
	return p.Parse(lmn)
}
