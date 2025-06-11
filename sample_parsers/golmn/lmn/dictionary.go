package lmn

import "maps"

func (lp *LmnParser) dictionary(skipFirst bool) (map[string]any, error) {
	if skipFirst { // skip first '(', when top-level Dictionary, it is false
		lp.idx++
	}
	lp.skip()

	var res = map[string]any{}
	var keys = map[string]struct{}{}
	var canSpread = true

	for lp.here() != ')' {
		if lp.here() == '.' && canSpread { // spread statement (..)
			lp.idx++
			if err := lp.consume('.'); err != nil {
				return nil, err
			}

			if ancVal, err := lp.getAnchorValue(); err != nil {
				return nil, err
			} else {
				switch ancVal.(type) {
				case map[string]any:
					maps.Copy(res, ancVal.(map[string]any))
				default:
					return nil, lp.err(mismatchAncTypeErr)
				}
			}
		} else {
			canSpread = false
			key, err := lp.key()

			if err != nil {
				return nil, err
			} else if _, exist := keys[key]; exist {
				return nil, lp.err(duplicatedKeyErr)
			}

			keys[key] = struct{}{}
			lp.skip()

			if lp.here() != ':' {
				// omit value with anchor
				if val, exist := lp.anchors[key]; exist {
					res[key] = val
					lp.skip()
				} else {
					return nil, lp.err(unexpectedTokenErr)
				}
			} else {
				lp.idx++

				if val, err := lp.value(); err != nil {
					return nil, err
				} else {
					res[key] = val
					lp.skip()
				}
			}
		}

		if here := lp.here(); here == ',' {
			lp.idx++
			lp.skip()
		} else if here == ')' {
			break
		} else {
			lp.err(unexpectedTokenErr)
		}
	}

	lp.idx++

	return res, nil
}

func (lp *LmnParser) topLevelDictionary() (map[string]any, error) {
	lp.buf = append(lp.buf, '\n') // to ignore comment
	lp.buf = append(lp.buf, ')')

	var res, err = lp.dictionary(false)

	lp.buf = lp.buf[:len(lp.buf)-2] // 버퍼 원상복구

	return res, err
}
