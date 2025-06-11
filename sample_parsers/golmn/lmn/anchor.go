package lmn

import "unicode/utf8"

func (lp *LmnParser) anchor() (string, error) {
	// 시작 허용 문자: 알파벳, _
	// 시작 비허용 문자: 숫자, -, .
	// 비어 있을 수 없음
	var res = []byte{}

	if here := lp.here(); '0' <= here && here <= '9' || here == '-' || here == '.' {
		return "", lp.err(identStartErr)
	}

	for here := lp.here(); 'a' <= here && here <= 'z' ||
		'A' <= here && here <= 'Z' ||
		'0' <= here && here <= '9' ||
		here == '_' || here == '-' || here == '.'; here = lp.here() {
		res = append(res, here)
		lp.idx++
	}

	if len(res) == 0 {
		return "", lp.err(emptyAnchorNameErr)
	}

	if utf8.Valid(res) {
		return string(res), nil
	} else {
		return "", lp.err(invalidEncodeErr)
	}
}

func (lp *LmnParser) getAnchorValue() (any, error) {
	anc, err := lp.anchor()

	if ancVal, exist := lp.anchors[anc]; exist && err == nil {
		return ancVal, nil
	} else {
		return nil, lp.err(failGetAncErr)
	}
}
