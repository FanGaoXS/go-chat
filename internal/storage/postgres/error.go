package postgres

import (
	"database/sql"
	stderr "errors"

	"fangaoxs.com/go-chat/internal/infras/errors"

	"github.com/lib/pq"
)

func wrapPGErrorf(err error, format string, args ...any) error {
	if stderr.Is(err, sql.ErrNoRows) {
		return errors.Newf(errors.NotFound, err, format, args...)
	}

	if e, ok := err.(*pq.Error); ok {
		switch e.Code.Name() {
		case "integrity_constraint_violation", "unique_violation":
			return errors.Newf(errors.AlreadyExists, err, format, args...)
		default:
			return errors.Newf(errors.Internal, err, format, args...)
		}
	}

	return errors.Newf(errors.Internal, err, format, args...)
}
