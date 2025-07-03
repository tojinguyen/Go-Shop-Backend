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
