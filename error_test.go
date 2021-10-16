// Copyright (c) 2021 Timo Savola. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package important_test

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/tsavola/important"
)

func Example(t *testing.T) {
	baseline := important.Unseen()
	err, seen := important.ErrorSeen(errors.New("forget me not"))

	fmt.Println(err)        // Observed: calls Error method.
	fmt.Printf("%v\n", err) // Observed: calls Error method.
	fmt.Printf("%s\n", err) // Observed (but incorrect usage).

	var linkError *os.LinkError

	errors.Unwrap(err)          // Observed.
	errors.Is(err, os.ErrExist) // Observed.
	errors.As(err, &linkError)  // Observed.
	os.IsExist(err)             // Not observed (IsExist is a legacy API).

	fmt.Errorf("%v", err) // Observed.
	fmt.Errorf("%w", err) // Observed.

	if !seen() {
		t.Fail()
	}
	if important.Unseen() != baseline {
		t.Fail()
	}
}

func TestSeen(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	if seen() {
		t.Fail()
	}

	err.Error()
	if !seen() {
		t.Fail()
	}

	err.Error()
	if !seen() {
		t.Fail()
	}
}

func TestPrint(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	fmt.Sprint(err)
	if !seen() {
		t.Fail()
	}
}

func TestFormatV(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	fmt.Sprintf("%v", err)
	if !seen() {
		t.Fail()
	}
}

func TestFormatS(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	fmt.Sprintf("%s", err)
	if !seen() { // ?
		t.Fail()
	}
}

func TestErrorV(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	fmt.Errorf("%v", err)
	if !seen() {
		t.Fail()
	}
}

func TestErrorW(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	fmt.Errorf("%w", err)
	if !seen() {
		t.Fail()
	}
}

func TestUnwrap(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	if errors.Unwrap(err) != io.EOF {
		t.Fail()
	}
	if !seen() {
		t.Fail()
	}
}

func TestIs(t *testing.T) {
	err, seen := important.ErrorSeen(io.EOF)
	if !errors.Is(err, io.EOF) {
		t.Fail()
	}
	if !seen() {
		t.Fail()
	}
}

type partyError struct{ expected int }

func (partyError) Error() string { return "party error" }

func TestAs(t *testing.T) {
	err, seen := important.ErrorSeen(&partyError{1999})

	var year *partyError
	if !errors.As(err, &year) {
		t.FailNow()
	}
	if year.expected != 1999 {
		t.Fail()
	}

	if !seen() {
		t.Fail()
	}
}

func TestUnseen(t *testing.T) {
	baseline := important.Unseen()

	err := important.Error(io.EOF)
	if n := important.Unseen() - baseline; n != 1 {
		t.Error(n)
	}

	err.Error()
	if n := important.Unseen() - baseline; n != 0 {
		t.Error(n)
	}

	important.Error(io.EOF)
	important.Error(io.EOF)
	if n := important.Unseen() - baseline; n != 2 {
		t.Error(n)
	}
}
