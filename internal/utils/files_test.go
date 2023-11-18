package utils

import (
    "testing"
)

func TestUpsertFolder(t *testing.T) {
    tests := []struct {
        name string
        input string
        errorExpected bool
    }{
        {
            name: "no trailing slash",
            input: "hello",
        },
        {
            name: "trailing slash",
            input: "hello/",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := UpsertFolder(test.input)
            if result != nil {
                t.Errorf("Error encountered on input %s: %s", test.input, result)
            }
        })
    }
}

func TestUpsertFile(t *testing.T) {
    tests := []struct {
        name string
        input string
        expected string
    }{
        {
            name: "no trailing slash",
            input: "hello",
        },
        {
            name: "trailing slash",
            input: "hello/",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := UpsertFile(test.input)
            if result != nil {
                t.Errorf("Error encountered on input %s: %s", test.input, result)
            }
        })
    }
}
