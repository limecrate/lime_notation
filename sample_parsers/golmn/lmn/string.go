package lmn

import (
	"unicode/utf8"
)

// Rust 문자열과 유사
// 단, 큰따옴표(") 대신 작은따옴표만 사용(') 작은따옴표 이스케이프 필요
// escapes map에 있는 문자 이스케이프 지원
// 유니코드 이스케이프 \u{AC00} 같이 중괄호 안에 최대 6자리 0x10ffff 이하의 코드 포인트 허용
// 아스키 헥사 이스케이프 \x21 같이 고정 2자리 16진수 표현 0x7f 이하 허용
// 문자열 내에 바로 개행 가능
// \이후 개행이 있다면 그 개행을 포함하여 연속된 공백 무시
// UTF-8만 지원
// Rust의 Raw String 같은 특수 문자열 지원 안함

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
			/* case 'x': // 헥사 아스키 코드 처리
			lp.idx++
			codePoint = 0

			for range 2 {
				if !isHexNum(lp.here()) {
					return "", lp.err(invalidEscErr)
				}
				codePoint = codePoint<<4 + rune(hexToInt(lp.here()))
				lp.idx++
			}

			if codePoint > 0x7f {
				return "", lp.err(invalidEscErr)
			}

			res = append(res, byte(codePoint)) */
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

func (lp *LmnParser) ident() (string, error) {
	if lp.here() == '\'' || lp.here() == '"' { // 문자열 형식의 키
		return lp.string()
	}
	return lp.anchor()
}
