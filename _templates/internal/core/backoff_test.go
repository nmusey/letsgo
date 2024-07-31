package core

import (
	"errors"
	"testing"
)

var failedAtempts = 0

var tests = []struct {
    attempts int
    fn func() error 
}{
    {1, func() error { return nil }},
    {2, func() error { 
        if failedAtempts < 1 {
            failedAtempts++
            return errors.New("Testing failed attempt")
        }

        return nil
    }},
}

func TestBlockingBackoff(t *testing.T) {
    for _, test := range tests {
        if err := BlockingBackoff(test.fn, test.attempts, 1); err != nil {
            t.Errorf("BlockingBackoff failed with %d attempts", test.attempts)
        }
    } 
}
