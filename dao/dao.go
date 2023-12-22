// DAOs, or Data Access Objects, are the mechanism used to talk to backing data stores.
package dao

import "errors"

var (
	NotFound error = errors.New("Record not found")
)
