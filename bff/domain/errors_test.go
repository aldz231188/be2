package domain

import "testing"

func TestValidationErrorsAddAndError(t *testing.T) {
	ve := NewValidationErrors()
	ve.Add("field1", "must not be empty")
	ve.Add("field2", "invalid")

	if !ve.HasErrors() {
		t.Fatal("expected errors to be present")
	}

	expected := "field1: must not be empty, field2: invalid"
	if ve.Error() != expected {
		t.Fatalf("unexpected error string: %q", ve.Error())
	}
}

func TestValidationErrorsMerge(t *testing.T) {
	first := NewValidationErrors()
	first.Add("first", "bad")

	second := NewValidationErrors()
	second.Add("second", "also bad")

	first.Merge(second)

	if len(first.Errors()) != 2 {
		t.Fatalf("expected 2 errors after merge, got %d", len(first.Errors()))
	}
}

func TestValidationErrorsNilSafety(t *testing.T) {
	var ve *ValidationErrors
	if ve.HasErrors() {
		t.Fatal("nil validation errors should not report errors")
	}

	// methods should not panic on nil receiver
	ve.Add("ignored", "ignored")
	ve.Merge(NewValidationErrors())
	if ve.Error() != "" {
		t.Fatal("expected empty string for nil ValidationErrors")
	}
}
