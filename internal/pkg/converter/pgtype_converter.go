package converter

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UUID conversions
func UUIDToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func PgUUIDToUUID(u pgtype.UUID) uuid.UUID {
	if u.Valid {
		return uuid.UUID(u.Bytes)
	}
	return uuid.Nil
}

func PgUUIDToString(u pgtype.UUID) string {
	if u.Valid {
		return u.String()
	}
	return ""
}

func StringToPgUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	u.Scan(s)
	return u
}

func NullPgUUID() pgtype.UUID {
	return pgtype.UUID{}
}

// Date conversions
func StringToPgDate(s string) pgtype.Date {
	var d pgtype.Date
	if s != "" {
		d.Time, _ = time.Parse("2006-01-02", s)
		d.Valid = true
	}
	return d
}

func PgDateToString(d pgtype.Date) string {
	if d.Valid {
		return d.Time.Format("2006-01-02")
	}
	return ""
}

// Timestamp conversions
func PgTimeToString(t pgtype.Timestamptz) string {
	if t.Valid {
		return t.Time.Format(time.RFC3339)
	}
	return ""
}

func NullPgTime() pgtype.Timestamptz {
	return pgtype.Timestamptz{}
}

func TimeToPgTime(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t, Valid: true}
}

// String to Timestamp conversion
func StringToPgTime(s string) pgtype.Timestamptz {
	var t pgtype.Timestamptz
	if s != "" {
		parsedTime, err := time.Parse(time.RFC3339, s)
		if err == nil {
			t.Time = parsedTime
			t.Valid = true
		}
	}
	return t
}

// Text conversions
func PgTextToStringPtr(t pgtype.Text) *string {
	if t.Valid {
		return &t.String
	}
	return nil
}

func StringToPgText(s *string) pgtype.Text {
	if s != nil {
		return pgtype.Text{String: *s, Valid: true}
	}
	return pgtype.Text{}
}

// Float conversions
func PgFloat8ToFloat64Ptr(f pgtype.Float8) *float64 {
	if f.Valid {
		return &f.Float64
	}
	return nil
}

func Float64ToPgFloat8(f *float64) pgtype.Float8 {
	if f != nil {
		return pgtype.Float8{Float64: *f, Valid: true}
	}
	return pgtype.Float8{}
}

// Add this to your converter package
func PgNumericToFloat64Ptr(numeric pgtype.Numeric) *float64 {
	if !numeric.Valid {
		return nil
	}

	float64Val, err := numeric.Float64Value()
	if err != nil {
		return nil
	}

	return PgFloat8ToFloat64Ptr(float64Val)
}

// Time pointer conversions
func PgTimeToTimePtr(t pgtype.Timestamptz) *time.Time {
	if t.Valid {
		return &t.Time
	}
	return nil
}

func TimePtrToPgTime(t *time.Time) pgtype.Timestamptz {
	if t != nil {
		return pgtype.Timestamptz{Time: *t, Valid: true}
	}
	return pgtype.Timestamptz{}
}

// Bool conversions
func BoolToPgBool(b bool) pgtype.Bool {
	return pgtype.Bool{Bool: b, Valid: true}
}

func PgBoolToBool(b pgtype.Bool) bool {
	if b.Valid {
		return b.Bool
	}
	return false
}
