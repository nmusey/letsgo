package utils

import (
    "testing"
)

func TestReplaceAllInString(t *testing.T) {
    tests := []struct {
        name string
        input string
        replacements map[string]string
        expected string
    }{
        {
            name: "no replacements",
            input: "hello world",
            replacements: map[string]string{},
            expected: "hello world",
        },
        {
            name: "one replacement",
            input: "hello world",
            replacements: map[string]string{
                "hello": "goodbye",
            },
            expected: "goodbye world",
        },
        {
            name: "multiple replacements",
            input: "hello world",
            replacements: map[string]string{
                "hello": "goodbye",
                "world": "universe",
            },
            expected: "goodbye universe",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := ReplaceAllInString(test.input, test.replacements)
            if result != test.expected {
                t.Errorf("expected %s, got %s", test.expected, result)
            }
        })
    }
}
