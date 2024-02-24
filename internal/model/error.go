package model

import "errors"

var (
	ErrThresholdIsZero  = errors.New("threshold is zero")
	ErrNotEnoughSamples = errors.New("not enough samples")
)
