package lmn

func (lp *LmnParser) list() ([]any, error) {
	lp.idx++
	lp.skip()

	var res = []any{}

	for lp.here() != ']' {
		if lp.here() == '.' { // spread statement (..)
			lp.idx++
			if err := lp.consume('.'); err != nil {
				return nil, err
			}

			if ancVal, err := lp.getAnchorValue(); err != nil {
				return nil, err
			} else {
				switch ancVal.(type) {
				case []any:
					res = append(res, ancVal.([]any)...)
				default:
					return nil, lp.err(mismatchAncTypeErr)
				}
			}

		} else if val, err := lp.value(); err != nil {
			return nil, err
		} else {
			res = append(res, val)
			lp.skip()
		}

		if here := lp.here(); here == ',' {
			lp.idx++
			lp.skip()
		} else if here == ']' {
			break
		} else {
			return nil, lp.err(unexpectedTokenErr)
		}
	}

	lp.idx++

	return res, nil
}
