package core

import (
    "time"
)

func BlockingBackoff(fn func() error, attempts int, sleep time.Duration) error {
    var err error
    for i := 0; i < attempts; i++ {
        if err := fn(); err == nil {
            return nil
        }

        time.Sleep(sleep)
    }

    return err
}
