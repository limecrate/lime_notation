package lmn

import "slices"

// var numbers = [...]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var whites = [...]byte{'\r', '\n', '\t', ' '}

func (lp LmnParser) here() byte {
	if lp.end() {
		return 0
	}
	return lp.buf[lp.idx]
}

func (lp LmnParser) end() bool {
	return lp.idx >= len(lp.buf)
}

func (lp LmnParser) notEnd() bool {
	return lp.idx < len(lp.buf)
}

func (lp *LmnParser) skipWhite() {
	for slices.Contains(whites[:], lp.here()) {
		lp.idx++
	}
}

func (lp *LmnParser) skipComment() {
	if lp.here() == '#' {
		lp.idx++

		for lp.notEnd() && lp.here() != '\n' {
			lp.idx++
		}

		if lp.notEnd() {
			lp.idx++
		}
	}
}

func (lp *LmnParser) skip() {
	for {
		idx := lp.idx
		lp.skipWhite()
		lp.skipComment()

		if idx == lp.idx {
			break
		}
	}
}
