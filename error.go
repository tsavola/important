// Copyright (c) 2021 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package important can be used to flag returned error values as important.
// Important errors need to be observed.  Observation means invocation of the
// the Unwrap method.  Unobserved errors cause no side-effects by default, but
// they may be checked in tests or tracked otherwise.
package important

import (
	"errors"
	"sync/atomic"
)

var unseen int64 // Atomic.

type errorType struct {
	err  error
	seen int64 // Atomic.
}

func (e *errorType) Error() string {
	return e.err.Error()
}

func (e *errorType) Unwrap() error {
	if atomic.AddInt64(&e.seen, 1) == 1 {
		atomic.AddInt64(&unseen, -1)
	}
	return e.err
}

// Error wraps the error, flagging it as important.
func Error(err error) error {
	err, _ = ErrorSeen(err)
	return err
}

// ErrorSeen wraps the error, flagging it as important.  The returned function
// may be called to see if the error has been observed.  It can be called many
// times, at any time.
func ErrorSeen(err error) (error, func() bool) {
	atomic.AddInt64(&unseen, 1)
	e := &errorType{err: err}
	f := func() bool { return atomic.LoadInt64(&e.seen) != 0 }
	return e, f
}

// Unwrap until the outermost important error, and return the error it wraps.
// If none is found, Unwrap returns nil.
//
// This can be used to flag a nested error as having been observed.
func Unwrap(err error) error {
	var e *errorType
	if errors.As(err, &e) {
		return e.Unwrap()
	}
	return nil
}

// Unseen error count since the start of the program.
func Unseen() int64 {
	return atomic.LoadInt64(&unseen)
}
