package mock

import (
	"encoding/json"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// ParametrizedBodyComparator can compare request and response with templates in values
// For example:
// {
//	 "name": "{name}"
// }
//
// and this mock will equal to any request with the same json structure
type ParametrizedBodyComparator struct{}

func (cmp ParametrizedBodyComparator) Compare(mock, request []byte) (bool, error) {

	var mockArr []interface{}
	var requestArr []interface{}

	err := json.Unmarshal(mock, &mockArr)
	err = json.Unmarshal(request, &requestArr)

	if err != nil {
		var mockObj map[string]interface{}
		var requestObj map[string]interface{}

		if err = json.Unmarshal(mock, &mockObj); err != nil {
			return false, err
		}

		if err = json.Unmarshal(request, &requestObj); err != nil {
			return false, err
		}

		return compareObject(mockObj, requestObj), nil
	}
	return compareAnyArrays(mockArr, requestArr), nil
}

func compareAnyArrays(mock, req []interface{}) bool {

	if len(mock) != len(req) {
		return false
	}

	for i := 0; i < len(mock); i++ {
		switch mt := mock[i].(type) {
		case map[string]interface{}:

			rv, ok := req[i].(map[string]interface{})
			if !ok {
				return false
			}

			result := make([]map[string]interface{}, len(mock))

			for i, m := range mock {
				result[i] = m.(map[string]interface{})
			}

			if !compareObject(mt, rv) {
				return false
			}
		default:
			if !reflect.DeepEqual(mock, req) {
				return false
			}
		}
	}

	return true
}

func compareObject(mock, request map[string]interface{}) bool {
	// Json is structure that can contains nested variables.
	// and two JSON can be equal only if they have equal number of keys
	// firstly check that keys count is equal

	if len(mock) != len(request) {
		return false
	}

	for key, value := range mock {
		reqVal, ok := request[key]

		// if there isn't the key in request then request is not equal to mock
		if !ok {
			return false
		}

		// check that mock value is string

		strValue, ok := value.(string)
		trimmed := strings.TrimSpace(strValue)

		if ok && len(trimmed) >= 2 {

			if trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}' {
				if !calculatePatternExpression(trimmed, reqVal) {
					return false
				}
				continue
			}
		}

		mockObj, ok := value.(map[string]interface{})

		if ok {
			reqObj, ok := reqVal.(map[string]interface{})

			if ok {
				// then both are objects
				return compareObject(mockObj, reqObj)
			} else {
				return false
			}
		}

		// it's not an object but in can be an array
		switch mt := value.(type) {
		case []interface{}:
			rv, ok := reqVal.([]interface{})
			if !ok {
				return false
			}
			if !compareAnyArrays(mt, rv) {
				return false
			}
		default:
			if !reflect.DeepEqual(value, reqVal) {
				return false
			}
		}
	}

	return true
}

func calculatePatternExpression(pattern string, requestValue interface{}) bool {

	patternLen := len(pattern) - 1
	patternExpression := strings.TrimSpace(pattern[1:patternLen])

	split := strings.Split(patternExpression, " ")

	if len(split) != 3 {
		return true
	}

	operation := split[1]

	right := split[2]

	switch operation {
	case "!=":
		return checkInequality(right, requestValue)
	case ">":
		switch vt := requestValue.(type) {
		case float64:
			val, err := strconv.ParseFloat(right, 64)
			if err != nil {
				log.Println("[ERR] cant convert right value to float64", right)
				return false
			}

			return vt > val
		default:
			log.Println("[ERR] incomparable type", vt)
			return false
		}
	case "<":
		switch vt := requestValue.(type) {
		case float64:
			val, err := strconv.ParseFloat(right, 64)
			if err != nil {
				log.Println("[ERR] cant convert right value to float64", right)
				return false
			}

			return vt < val
		default:
			log.Println("[ERR] incomparable type", vt)
			return false
		}
	case "<=":
		switch vt := requestValue.(type) {
		case float64:
			val, err := strconv.ParseFloat(right, 64)
			if err != nil {
				log.Println("[ERR] cant convert right value to float64", right)
				return false
			}

			return vt <= val
		default:
			log.Println("[ERR] incomparable type", vt)
			return false
		}
	case ">=":
		switch vt := requestValue.(type) {
		case float64:
			val, err := strconv.ParseFloat(right, 64)
			if err != nil {
				log.Println("[ERR] cant convert right value to float64", right)
				return false
			}

			return vt >= val
		default:
			log.Println("[ERR] incomparable type", vt)
			return false
		}
	default:
		log.Println("[ERR] unsupported operation", operation)
		return true
	}
}

func checkInequality(l string, r interface{}) bool {
	switch vt := r.(type) {
	case int:
		val, err := strconv.Atoi(l)
		if err != nil {
			log.Println("[ERR] cant convert right value to int", l)
			return false
		}
		return val != vt
	case float64:
		val, err := strconv.ParseFloat(l, 64)
		if err != nil {
			log.Println("[ERR] cant convert right value to float64", l)
			return false
		}

		return val != vt
	case string:
		return l != vt
	default:
		log.Println("[ERR] incomparable type", vt)
		return false
	}
}
