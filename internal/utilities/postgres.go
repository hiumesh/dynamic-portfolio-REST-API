package utilities

import (
	"errors"
	"strconv"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresError struct {
	Code           string `json:"code"`
	HttpStatusCode int    `json:"-"`
	Message        string `json:"message"`
	Hint           string `json:"hint,omitempty"`
	Detail         string `json:"detail,omitempty"`
}

func NewPostgresError(err error) *PostgresError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && isPubliclyAccessiblePostgresError(pgErr.Code) {
		return &PostgresError{
			Code:           pgErr.Code,
			HttpStatusCode: getHttpStatusCodeFromPostgresErrorCode(pgErr.Code),
			Message:        pgErr.Message,
			Detail:         pgErr.Detail,
			Hint:           pgErr.Hint,
		}
	}

	return nil
}

func isPubliclyAccessiblePostgresError(code string) bool {
	if len(code) != 5 {
		return false
	}

	// default response
	return getHttpStatusCodeFromPostgresErrorCode(code) != 0
}

func getHttpStatusCodeFromPostgresErrorCode(code string) int {
	if code == pgerrcode.RaiseException ||
		code == pgerrcode.IntegrityConstraintViolation ||
		code == pgerrcode.RestrictViolation ||
		code == pgerrcode.NotNullViolation ||
		code == pgerrcode.ForeignKeyViolation ||
		code == pgerrcode.UniqueViolation ||
		code == pgerrcode.CheckViolation ||
		code == pgerrcode.ExclusionViolation {
		return 500
	}

	// Use custom HTTP status code if Postgres error was triggered with `PTXXX`
	// code. This is consistent with PostgREST's behaviour as well.
	if strings.HasPrefix(code, "PT") {
		if httpStatusCode, err := strconv.ParseInt(code[2:], 10, 0); err == nil {
			return int(httpStatusCode)
		}
	}

	return 0
}
