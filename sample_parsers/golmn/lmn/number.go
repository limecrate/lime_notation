package lmn

import "strconv"

// 정수
// 0x, 0o, 0b 지원
// Zero-Padding 제한 없음
// _ 구분의 경우 진수 표시, 지수 뒤 혹은 지수의 부호 뒤에 바로 오지만 않으면 제한 없음
// 진수 표시자, 지수(e, E) 다음에는 숫자가 있어야함 (아무것도 없으면 안됨)
// 잘못된 예: _342, -0b_324, 34E+_342, 0x, 3.14E_
// 실수 소수점, 지수 표기 방식 지원
// 소수점 뒤에 아무것도 없어도 됨 0., 3.__ 같은 형식 지원

func (lp *LmnParser) number(sign int) (any, error) {
	var isInt = true
	var resInt = 0
	var resNums = []byte{}

	if lp.here() == '0' {
		lp.idx++

		if next := lp.here(); next == 'x' { // base 16
			return lp.intWithBase(4, sign)
		} else if next == 'o' { // base 8
			return lp.intWithBase(3, sign)
		} else if next == 'b' { // base 2
			return lp.intWithBase(1, sign)
		}
	}

	for here := lp.here(); '0' <= here && here <= '9' || here == '_'; here = lp.here() {
		if here != '_' {
			resInt = resInt*10 + int(here-'0')
		}
		lp.idx++
	}

	if lp.here() == '.' {
		isInt = false
		lp.idx++
		resNums = strconv.AppendInt(resNums, int64(resInt), 10)
		resNums = append(resNums, '.')

		for here := lp.here(); isDecNum(here) || here == '_'; here = lp.here() {
			if here != '_' {
				resNums = append(resNums, here)
			}
			lp.idx++
		}
	}

	if lp.here() == 'e' || lp.here() == 'E' {
		if len(resNums) == 0 { // resNums 준비
			resNums = strconv.AppendInt(resNums, int64(resInt), 10)
		}

		isInt = false
		resNums = append(resNums, 'E')
		lp.idx++

		if here := lp.here(); here == '+' || here == '-' {
			resNums = append(resNums, here)
			lp.idx++
		}

		if here := lp.here(); !isDecNum(here) {
			return 0.0, lp.err(expectNumErr)
		}

		for here := lp.here(); isDecNum(here) || here == '_'; here = lp.here() {
			if here != '_' {
				resNums = append(resNums, here)
			}
			lp.idx++
		}
	}

	if isInt {
		return int(sign) * resInt, nil
	} else {
		if resNum, err := strconv.ParseFloat(string(resNums), 64); err != nil {
			return 0.0, err
		} else {
			return float64(sign) * resNum, nil
		}
	}
}

func (lp *LmnParser) intWithBase(baseExp int, sign int) (int, error) {
	lp.idx++

	var res int = 0

	if baseExp == 1 && !isBinNum(lp.here()) ||
		baseExp == 3 && !isOctNum(lp.here()) ||
		baseExp == 4 && !isHexNum(lp.here()) {
		return 0, lp.err(expectNumErr)
	}

	switch baseExp {
	case 1, 3:
		var end byte

		if baseExp == 1 {
			end = '1'
		} else {
			end = '7'
		}

		for here := lp.here(); '0' <= here && here <= end || here == '_'; here = lp.here() {
			if here != '_' {
				res = res<<baseExp + int(here-'0')
			}
			lp.idx++
		}
	case 4:
		for here := lp.here(); ; here = lp.here() {
			if '0' <= here && here <= '9' {
				res = res<<4 + int(here) - '0'
			} else if 'a' <= here && here <= 'f' {
				res = res<<4 + int(here) - 'a' + 10
			} else if 'A' <= here && here <= 'F' {
				res = res<<4 + int(here) - 'A' + 10
			} else if here != '_' {
				break
			}
			lp.idx++
		}
	}

	return int(sign) * res, nil
}

func isHexNum(n byte) bool {
	return '0' <= n && n <= '9' ||
		'a' <= n && n <= 'f' ||
		'A' <= n && n <= 'F'
}

func isOctNum(n byte) bool {
	return '0' <= n && n < '7'
}

func isDecNum(n byte) bool {
	return '0' <= n && n < '9'
}

func isBinNum(n byte) bool {
	return '0' <= n && n < '1'
}

func hexToInt(h byte) int {
	if '0' <= h && h <= '9' {
		return int(h) - '0'
	} else if 'a' <= h && h <= 'f' {
		return int(h) - 'a' + 10
	} else if 'A' <= h && h <= 'F' {
		return int(h) - 'A' + 10
	} else {
		return -1
	}
}
