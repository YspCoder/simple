package builder

import (
	"github.com/YspCoder/simple/common/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// appendNotNull appends the provided key and value to the map if the value is not nil.
func appendNotNull(m bson.M, key string, val interface{}) {
	if !utils.IsNil(val) {
		m[key] = val
	}
}
