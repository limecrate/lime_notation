package lmn

import (
	"errors"
	"fmt"
)

var expectValueErr = errors.New("Value Expected")
var expectNumErr = errors.New("Number Expected")
var unexpectedTokenErr = errors.New("Unexpect Token")
var invalidEscErr = errors.New("Invalid Escape Sequence")
var invalidEncodeErr = errors.New("Invalid String Sequence")
var emptyAnchorNameErr = errors.New("Empty Anchor Name")
var duplicatedKeyErr = errors.New("This Key is Duplicated")
var duplicatedCapErr = errors.New("This Capture is Duplicated")
var identStartErr = errors.New("Identifier Cannot Start with This Letter")
var failGetAncErr = errors.New("Fail to Get Anchor's Value")
var mismatchAncTypeErr = errors.New("Mismatched Anchor Type")
var temErr = errors.New("Temporary Error")

func (lp LmnParser) err(err error) error {
	return fmt.Errorf("Unexpected %d at %d: %w", lp.here(), lp.idx, err)
}

func (lp *LmnParser) consume(c byte) error {
	if lp.here() != c {
		return lp.err(unexpectedTokenErr)
	}
	lp.idx++
	return nil
}
