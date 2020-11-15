package mgmrel

import (
	"fmt"
	"github.com/kamva/gutil"
	"github.com/kamva/mgm/v3"
	"reflect"
)

// foreignKeyName gets the Model and returns foreignKey field name.
// e.g bet the "Book" model and returns "book_id".
func foreignKeyName(m mgm.Model) string {
	name := reflect.TypeOf(m).Elem().Name()
	return fmt.Sprintf("%s_id", gutil.ToSnakeCase(name))
}
