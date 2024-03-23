package postgres

import "fmt"

func newListError(model string, err error) error {
	return fmt.Errorf("failed to query table %s: %w", model, err)
}

func newInsertError(model string, err error) error {
	return fmt.Errorf("failed to insert table %s: %w", model, err)
}

func newNoRowsError(model string, err error) error {
	return fmt.Errorf("failed to return rows for table %s: %w", model, err)
}
