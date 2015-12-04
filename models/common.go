package models

import (
	"fmt"
)

func gtkErr(myfunc, gtkfunc string, err error) error {
	return fmt.Errorf("%s: Function %s returned an error: %v",
		myfunc, gtkfunc, err)
}
