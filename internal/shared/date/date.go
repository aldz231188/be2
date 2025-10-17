package date

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const dateLayout = "2006-01-02"

// DateOnly — дата без времени (нормализуем к полуночи UTC)
type DateOnly struct{ time.Time }

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	// ожидаем "YYYY-MM-DD" или null
	if string(b) == "null" {
		d.Time = time.Time{}
		return nil
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("date must be string %q: %w", dateLayout, err)
	}
	s = strings.TrimSpace(s)
	t, err := time.ParseInLocation(dateLayout, s, time.UTC)
	if err != nil {
		return fmt.Errorf("bad date %q, want %s: %w", s, dateLayout, err)
	}
	d.Time = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}

// Полезно, если пишешь в БД как DATE
func (d DateOnly) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time, nil
}

func (d *DateOnly) Scan(src any) error {
	switch v := src.(type) {
	case time.Time:
		d.Time = time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, time.UTC)
		return nil
	case string:
		t, err := time.ParseInLocation(dateLayout, v, time.UTC)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	case []byte:
		t, err := time.ParseInLocation(dateLayout, string(v), time.UTC)
		if err != nil {
			return err
		}
		d.Time = t
		return nil
	default:
		return fmt.Errorf("unsupported Scan type %T", src)
	}
}
