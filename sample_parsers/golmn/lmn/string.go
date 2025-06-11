package lmn

import (
	"unicode/utf8"
)

// 작은따옴표(')와 큰따옴표(") 사용
// escapes map에 있는 문자 이스케이프 지원
// 유니코드 이스케이프 \u{AC00} 같이 중괄호 안에 최대 6자리 0x10ffff 이하의 코드 포인트 허용
// 문자열 내에 바로 개행 가능
// \이후 개행이 있다면 그 개행을 포함하여 연속된 공백 무시 (string continuation)
// UTF-8 지원

var escapes = map[byte]byte{
	'\'': '\'',
	'\\': '\\',
	'"':  '"',
	'n':  '\n',
	'r':  '\r',
	't':  '\t',
}

func (lp *LmnParser) string() (string, error) {
	var res = []byte{}
	var codePoint rune = 0
	var uniBuf = make([]byte, 4)
	var quote = lp.here()

	lp.idx++

	for lp.here() != quote && !lp.end() {
		if here := lp.here(); here == '\\' {
			lp.idx++

			switch here = lp.here(); here {
			case 'u': // 유니코드 처리
				lp.idx++
				codePoint = 0

				if err := lp.consume('{'); err != nil {
					return "", err
				}

				if !isHexNum(lp.here()) {
					return "", lp.err(expectNumErr)
				}

				for i := 0; i < 6 && isHexNum(lp.here()); i++ {
					codePoint = codePoint<<4 + rune(hexToInt(lp.here()))
					lp.idx++
				}

				if err := lp.consume('}'); err != nil {
					return "", err
				}

				if codePoint > 0x10ffff {
					return "", lp.err(invalidEscErr)
				}

				n := utf8.EncodeRune(uniBuf, codePoint)
				res = append(res, uniBuf[:n]...)
			case '(': // 문자열 보간 처리, capture만 가능
				lp.idx++
				val, err := lp.getAnchorValue()

				if err != nil {
					return "", lp.err(failGetAncErr)
				}

				switch val.(type) {
				case string:
					res = append(res, []byte(val.(string))...)
					if err = lp.consume(')'); err != nil {
						return "", lp.err(unexpectedTokenErr)
					}
				default:
					return "", lp.err(mismatchAncTypeErr)
				}

			case '\\', '\'', '"', 'n', 'r', 't': // 이스케이프 처리
				res = append(res, escapes[here])
				lp.idx++
			case '\n': // string continuation
				lp.skipWhite()
			default:
				return "", lp.err(invalidEscErr)
			}
		} else {
			res = append(res, here)
			lp.idx++
		}
	}

	if err := lp.consume(quote); err != nil {
		return "", err
	}

	if utf8.Valid(res) {
		return string(res), nil
	} else {
		return "", lp.err(invalidEncodeErr)
	}
}

func (lp *LmnParser) key() (string, error) {
	if lp.here() == '\'' || lp.here() == '"' { // 문자열 형식 키
		return lp.string()
	}
	return lp.anchor() // 앵커 형식 키
}
