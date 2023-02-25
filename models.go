package requests

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Float float64

func (f *Float) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	var fv float64
	switch v := i.(type) {
	case int:
		fv = float64(v)
	case string:
		if v == "" {
			break
		}
		sf, err := strconv.ParseFloat(strings.ReplaceAll(v, ",", ""), 64)
		if err != nil {
			return err
		}
		fv = sf
	case float64:
		fv = v
	default:
		return fmt.Errorf("float: invalid type, %T", i)
	}

	*f = Float(fv)
	return nil
}
