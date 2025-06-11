package lmn

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

func LmnParse(lmn string) (any, error) {
	p := NewLmn()
	return p.Parse(lmn)
}
