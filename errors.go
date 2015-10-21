package sweet

import "fmt"

type ErrNotFound struct {
	Selection string
}

func (err ErrNotFound) Error() string {
	return fmt.Sprintf("Could not match any elements with selector %v", err.Selection)
}

type ErrEmptyForm struct {
	Selection string
}

func (err ErrEmptyForm) Error() string {
	return fmt.Sprintf("No input text to be posted within form %v", err.Selection)
}
