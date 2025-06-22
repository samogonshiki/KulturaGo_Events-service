package dto

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"unicode"
)

type FlexTime time.Time

func (ft FlexTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ft)
	b, err := json.Marshal(t.Format(time.RFC3339))
	return b, err
}

func (ft *FlexTime) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if len(s) >= 3 {
		suf := s[len(s)-3:]
		if (suf[0] == '+' || suf[0] == '-') &&
			unicode.IsDigit(rune(suf[1])) &&
			unicode.IsDigit(rune(suf[2])) &&
			!strings.Contains(suf, ":") {
			s += ":00"
		}
	}
	if strings.HasSuffix(s, "Z07:00") {
		s = strings.Replace(s, "Z07:00", "Z", 1)
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("parsing time %q: %w", s, err)
	}
	*ft = FlexTime(t)
	return nil
}

func (ft FlexTime) Value() (driver.Value, error) {
	return time.Time(ft), nil
}
