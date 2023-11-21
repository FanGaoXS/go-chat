package postgres

import "errors"

func wrapPGErrorf(err error) error {
	return errors.Newf
}
