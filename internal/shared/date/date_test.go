package date

import (
	"encoding/json"
	"testing"
	"time"
)

func TestDateOnlyJSONRoundtrip(t *testing.T) {
	input := []byte(`"2024-03-15"`)
	var d DateOnly
	if err := json.Unmarshal(input, &d); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	expected := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if !d.Time.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, d.Time)
	}

	out, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	if string(out) != "\"2024-03-15\"" {
		t.Fatalf("unexpected marshal output: %s", out)
	}
}

func TestDateOnlyZeroValueMarshalling(t *testing.T) {
	var d DateOnly
	out, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if string(out) != "null" {
		t.Fatalf("expected null for zero value, got %s", out)
	}
}

func TestDateOnlyInvalidJSON(t *testing.T) {
	var d DateOnly
	if err := json.Unmarshal([]byte(`"2024-13-40"`), &d); err == nil {
		t.Fatal("expected error for invalid date, got nil")
	}
}

func TestDateOnlyValueAndScan(t *testing.T) {
	expected := time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC)
	d := DateOnly{Time: expected}

	val, err := d.Value()
	if err != nil {
		t.Fatalf("value failed: %v", err)
	}
	if !val.(time.Time).Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, val)
	}

	var scanned DateOnly
	if err := scanned.Scan(expected); err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if !scanned.Time.Equal(expected) {
		t.Fatalf("expected %v after scan, got %v", expected, scanned.Time)
	}

	if err := scanned.Scan([]byte("2023-12-05")); err != nil {
		t.Fatalf("scan from bytes failed: %v", err)
	}
	if !scanned.Time.Equal(expected) {
		t.Fatalf("expected %v after byte scan, got %v", expected, scanned.Time)
	}
}

func TestDateOnlyScanUnsupported(t *testing.T) {
	var d DateOnly
	if err := d.Scan(123); err == nil {
		t.Fatal("expected error for unsupported scan type")
	}
}
