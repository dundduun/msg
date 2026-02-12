package werr

import "fmt"

func WrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
