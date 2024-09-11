package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var allowedBooleanMap = map[string]bool{
	"1":     true,
	"t":     true,
	"T":     true,
	"true":  true,
	"True":  true,
	"TRUE":  true,
	"0":     false,
	"f":     false,
	"F":     false,
	"false": false,
	"False": false,
	"FALSE": false,
}

func main() {
	jsonData := `
	{
  "number_1": {
    "N": "1.50"
  },
  "string_1": {
    "S": "784498 "
  },
  "string_2": {
    "S": "2014-07-16T20:55:46Z"
  },
  "map_1": {
    "M": {
      "bool_1": {
        "BOOL": "truthy"
      },
      "null_1": {
        "NULL ": "true"
      },
      "list_1": {
        "L": [
          {
            "S": ""
          },
          {
            "N": "011"
          },
          {
            "N": "5215s"
          },
          {
            "BOOL": "f"
          },
          {
            "NULL": "0"
          }
        ]
      }
    }
  },
  "list_2": {
    "L": "noop"
  },
  "list_3": {
    "L": [
      "noop"
    ]
  },
  "": {
    "S": "noop"
  }
}`

	result := make([]map[string]interface{}, 0, 1)
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("Cannot unmarshall the data")
	}

	obj := make(map[string]interface{})
	for key, val := range data {
		if trim(key) != "" {
			fillup(key, val, obj)
		}
	}

	result = append(result, obj)
	marshalled, err := json.MarshalIndent(result, "", "  ")
	if err == nil {
		fmt.Println(string(marshalled))
	}
}

func fillup(key string, val interface{}, obj map[string]interface{}) {
	typedVal := val.(map[string]interface{})

	if mapData, ok := typedVal["M"]; ok && mapData != nil {

		sanitizedBool := mapData.(map[string]interface{})
		newObj := make(map[string]interface{})
		for key, val := range sanitizedBool {
			// Remove keys with empty string
			if trim(key) != "" {
				fillup(key, val, newObj)
			}
		}
		if len(newObj) != 0 {
			obj[key] = newObj
		}

	} else if listData, ok := typedVal["L"]; ok && listData != nil {

		listKind := reflect.ValueOf(listData).Kind()
		if listKind == reflect.Array || listKind == reflect.Slice {
			sanitizedList := listData.([]interface{})
			newList := make([]interface{}, 0, 1)
			for _, val := range sanitizedList {
				listObjKind := reflect.ValueOf(val).Kind()
				if listObjKind == reflect.Map {
					scopedVal := val.(map[string]interface{})
					sanitizedScopedVal := sanitizeKeys(scopedVal)
					listValue, err := nth_level_extractor(sanitizedScopedVal)
					if err == nil {
						newList = append(newList, listValue)
					}
				}
			}
			if len(newList) != 0 {
				obj[key] = newList
			}
		}

	} else {
		sanitizedTypedVal := sanitizeKeys(typedVal)
		granularVal, err := nth_level_extractor(sanitizedTypedVal)
		if err == nil {
			obj[key] = granularVal
		}
	}

}

func nth_level_extractor(typedVal map[string]interface{}) (interface{}, error) {
	if num, ok := typedVal["N"]; ok && num != nil && num != "" {
		if numberVal, err := strconv.ParseInt(num.(string), 10, 64); err == nil {
			return numberVal, nil
		}
		if numberVal, err := strconv.ParseFloat(num.(string), 62); err == nil {
			return numberVal, nil
		}

	}

	if num, ok := typedVal["S"]; ok && num != nil && num != "" {

		parsedTime, err := time.Parse(time.RFC3339, num.(string))
		if err == nil {
			return parsedTime.Unix(), nil
		} else if trim(num) != "" {
			return trim(num), nil
		}

	}

	if boolData, ok := typedVal["BOOL"]; ok && boolData != nil && boolData != "" {

		sanitizedBool := boolData.(string)
		if trim(sanitizedBool) != "" {
			if boolVal, ok := allowedBooleanMap[sanitizedBool]; ok {
				return boolVal, nil
			}
		}

	}
	if nullData, ok := typedVal["NULL"]; ok && nullData != nil && nullData != "" {

		sanitizedNull := nullData.(string)
		if trim(sanitizedNull) != "" {
			if nullVal, ok := allowedBooleanMap[trim(sanitizedNull)]; ok {
				if nullVal {
					return nil, nil
				}
			}
		}

	}
	return nil, errors.New("Unmapped data")
}

func trim(val interface{}) string {
	return strings.TrimSpace(val.(string))
}

func sanitizeKeys(val map[string]interface{}) map[string]interface{} {
	response := make(map[string]interface{})
	for k, v := range val {
		response[trim(k)] = v
	}
	return response
}

