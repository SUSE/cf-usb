package strfmt

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func init() {
	dt := DateTime{}
	Default.Add("datetime", &dt, IsDateTime)
}

// IsDateTime returns true when the string is a valid date-time
func IsDateTime(str string) bool {
	s := strings.Split(strings.ToLower(str), "t")
	if !IsDate(s[0]) {
		return false
	}

	matches := rxDateTime.FindAllStringSubmatch(s[1], -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		return false
	}
	m := matches[0]
	res := m[1] <= "23" && m[2] <= "59" && m[3] <= "59"
	return res
}

const (
	// RFC3339Millis represents a ISO8601 format to millis instead of to nanos
	RFC3339Millis = "2006-01-02T15:04:05.000Z07:00"
	// DateTimePattern pattern to match for the date-time format from http://tools.ietf.org/html/rfc3339#section-5.6
	DateTimePattern = `^([0-9]{2}):([0-9]{2}):([0-9]{2})(.[0-9]+)?(z|([+-][0-9]{2}:[0-9]{2}))$`
)

var (
	dateTimeFormats = []string{RFC3339Millis, time.RFC3339, time.RFC3339Nano}
	rxDateTime      = regexp.MustCompile(DateTimePattern)
)

// ParseDateTime parses a string that represents an ISO8601 time or a unix epoch
func ParseDateTime(data string) (DateTime, error) {
	if data == "" {
		return DateTime{Time: time.Unix(0, 0).UTC()}, nil
	}
	var lastError error
	for _, layout := range dateTimeFormats {
		dd, err := time.Parse(layout, data)
		if err != nil {
			lastError = err
			continue
		}
		lastError = nil
		return DateTime{dd}, nil
	}
	return DateTime{}, lastError
}

// DateTime is a time but it serializes to ISO8601 format with millis
// It knows how to read 3 different variations of a RFC3339 date time.
// Most API's we encounter want either millisecond or second precision times. This just tries to make it worry-free.
//
// swagger:strfmt date-time
type DateTime struct {
	time.Time
}

func (t DateTime) String() string {
	return t.Format(RFC3339Millis)
}

// MarshalText implements the text marshaller interface
func (t DateTime) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

// UnmarshalText implements the text unmarshaller interface
func (t *DateTime) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}
	tt, err := ParseDateTime(string(text))
	if err != nil {
		return err
	}
	*t = tt
	return nil
}

// Scan scans a DateTime value from database driver type.
func (d *DateTime) Scan(raw interface{}) error {
	// TODO: case int64: and case float64: ?
	switch v := raw.(type) {
	case []byte:
		return d.UnmarshalText(v)
	case string:
		return d.UnmarshalText([]byte(v))
	case time.Time:
		*d = DateTime{v}
	case nil:
		*d = DateTime{}
	default:
		return fmt.Errorf("cannot sql.Scan() strfmt.DateTime from: %#v", v)
	}

	return nil
}

// Value converts DateTime to a primitive value ready to written to a database.
func (d DateTime) Value() (driver.Value, error) {
	return driver.Value(d.Time), nil
}
