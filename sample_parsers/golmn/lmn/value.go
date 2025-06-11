package lmn

import "math"

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

	// ~ 앵커 처리
	lp.skip()

	if lp.here() == '~' {
		lp.idx++
		lp.skip()

		if anc, err := lp.anchor(); err != nil {
			return nil, err
		} else {
			if _, exist := lp.anchors[anc]; exist {
				return nil, lp.err(duplicatedCapErr)
			} else {
				lp.anchors[anc] = val
			}
		}
	}

	return val, err
}
