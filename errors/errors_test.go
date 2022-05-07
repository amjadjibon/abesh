package errors

import (
	"errors"
	"testing"
)

var ErrTest = NewWithAllInfo(1, "T", "failed", map[string]string{}, true)

func TestError_Error(t *testing.T) {
	var err = New(1, "T", "failed", map[string]string{})
	if err.Error() != "T: failed" {
		t.Errorf("Error() = %s, want %s", err.Error(), "T: failed")
	}

	var err1 = NewInternalErrorWithCause(errors.New("cause"), "failed", nil, "1")
	if err1.Error() != "internal_service.1: failed: cause" {
		t.Errorf("Error() = %s, want %s", err1.Error(), "internal_service.1: failed: cause")
	}
}

func TestError_Unwrap(t *testing.T) {
	var err = New(1, "T", "failed", map[string]string{})
	if err.Unwrap() != nil {
		t.Errorf("Unwrap() = %s, want %v", err.Unwrap(), nil)
	}

	var err1 = NewInternalErrorWithCause(errors.New("cause"), "failed", nil, "1")
	if err1.Unwrap().Error() != errors.New("cause").Error() {
		t.Errorf("Unwrap() = %s, want %s", err1.Unwrap(), "cause")
	}
}

func TestError_legacyErrString(t *testing.T) {
	var err = New(1, "T", "failed", map[string]string{})
	if err.legacyErrString() != "T: failed" {
		t.Errorf("legacyErrString() = %s, want %s", err.legacyErrString(), "T: failed")
	}
}

func TestError_StackTrace(t *testing.T) {
	var err = New(1, "T", "failed", map[string]string{})
	for _, f := range err.StackTrace() {
		if f == 0 {
			t.Errorf("StackTrace() = %v, want %v", f, nil)
		}
	}
}

func TestError_VerboseString(t *testing.T) {
	var err = New(1, "T", "failed", map[string]string{})
	if err.VerboseString()[0] != 'T' {
		t.Errorf("VerboseString() = %s, want %s", err.VerboseString(), "")
	}
}

func TestError_PrefixMatches(t *testing.T) {
	var err = New(1, "T", "failed", map[string]string{})
	if err.PrefixMatches("T") != true {
		t.Errorf("PrefixMatches() = %v, want %v", err.PrefixMatches("T"), true)
	}

	var err1 = New(1, "T", "failed", map[string]string{})
	if err1.PrefixMatches("T1") != false {
		t.Errorf("PrefixMatches() = %v, want %v", err1.PrefixMatches("T1"), false)
	}
}

func TestPrefixMatches(t *testing.T) {
	var err = New(1, "t", "failed", map[string]string{})
	if PrefixMatches(err, "t") != true {
		t.Errorf("PrefixMatches() = %v, want %v", PrefixMatches(err, "t"), true)
	}

	var err1 = errors.New("t1: failed")
	if PrefixMatches(err1, "t2") != false {
		t.Errorf("PrefixMatches() = %v, want %v", PrefixMatches(err1, "t2"), false)
	}

}

func TestError_ProtoErrorWithStack(t *testing.T) {
	t.Logf("%+v\n", ErrTest.StackString())
	t.Logf("%+v\n", NewInternalErrorWithCause(errors.New("test"), "just test", nil, "missing").StackString())

	t.Logf("%+v\n", NewFromError(ErrTest).StackString())
}

func TestError_IsRetryable(t *testing.T) {
	if !IsRetryable(ErrTest) {
		t.Error("error should be retryable")
	}
}

func Test_IsPrefixMatches(t *testing.T) {
	if !IsPrefixMatches(ErrTest, "T") {
		t.Error("prefix should match with T")
	}
}
