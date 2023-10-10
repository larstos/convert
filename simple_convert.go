package convert

import (
	"encoding/json"
	"math"
	"strconv"
)

func MustString(v any, defaultval ...string) string {
	val, ok := TryString(v)
	if ok {
		return val
	}
	if len(defaultval) > 0 {
		return defaultval[0]
	}
	return ""
}

func TryString(v any) (string, bool) {
	switch tv := v.(type) {
	case string:
		return tv, true
	case []byte:
		return string(tv), true
	case int64:
		return strconv.FormatInt(int64(tv), 10), true
	case uint64:
		return strconv.FormatUint(uint64(tv), 10), true
	case int32:
		return strconv.FormatInt(int64(tv), 10), true
	case uint32:
		return strconv.FormatUint(uint64(tv), 10), true
	case int16:
		return strconv.FormatInt(int64(tv), 10), true
	case uint16:
		return strconv.FormatUint(uint64(tv), 10), true
	case int8:
		return strconv.FormatInt(int64(tv), 10), true
	case uint8:
		return strconv.FormatUint(uint64(tv), 10), true
	case float32:
		return strconv.FormatFloat(float64(tv), 'f', -1, 64), true
	case float64:
		return strconv.FormatFloat(float64(tv), 'f', -1, 64), true
	case int:
		return strconv.Itoa(int(tv)), true
	case json.Number:
		return tv.String(), true
	case bool:
		if tv {
			return "true", true
		} else {
			return "false", true
		}
	}
	return "", false
}

func MustInt64(v any, defaultval ...int64) int64 {
	var defaultValue int64 = 0
	if len(defaultval) > 0 {
		defaultValue = defaultval[0]
	}
	if v == nil {
		return defaultValue
	}
	switch tv := v.(type) {
	case float32:
		if tv > float32(math.MaxInt64) {
			return defaultValue
		}
		return int64(tv)
	case float64:
		if tv > float64(math.MaxInt64) {
			return defaultValue
		}
		return int64(tv)
	}
	val, ok := TryInt64(v)
	if ok {
		return val
	}
	return defaultValue
}

func TryInt64(v any) (int64, bool) {
	if v == nil {
		return -1, false
	}
	switch tv := v.(type) {
	case []byte:
		res, err := strconv.ParseInt(string(tv), 10, 0)
		if err != nil {
			return -1, false
		}
		return res, true
	case string:
		res, err := strconv.ParseInt(tv, 10, 0)
		if err != nil {
			return -1, false
		}
		return res, true
	case int64:
		return tv, true
	case uint64:
		if tv > uint64(math.MaxInt64) {
			return -1, false
		}
		return int64(tv), true
	case int32:
		return int64(tv), true
	case uint32:
		return int64(tv), true
	case int:
		return int64(tv), true
	case int16:
		return int64(tv), true
	case uint16:
		return int64(tv), true
	case int8:
		return int64(tv), true
	case uint8:
		return int64(tv), true
	case json.Number:
		val, err := tv.Int64()
		if err == nil {
			return val, true
		}
	}
	return -1, false
}

func MustFloat64(v any, defaultval ...float64) float64 {
	val, ok := TryFloat64(v)
	if ok {
		return val
	}
	if len(defaultval) > 0 {
		return defaultval[0]
	}
	return 0
}

func TryFloat64(v any) (float64, bool) {
	if v == nil {
		return 0, false
	}
	switch tv := v.(type) {
	case []byte:
		res, err := strconv.ParseFloat(string(tv), 0)
		if err != nil {
			return 0, false
		}
		return res, true
	case string:
		res, err := strconv.ParseFloat(tv, 0)
		if err != nil {
			return 0, false
		}
		return res, true
	case int64:
		return float64(tv), true
	case uint64:
		return float64(tv), true
	case int32:
		return float64(tv), true
	case uint32:
		return float64(tv), true
	case int:
		return float64(tv), true
	case int16:
		return float64(tv), true
	case uint16:
		return float64(tv), true
	case int8:
		return float64(tv), true
	case uint8:
		return float64(tv), true
	case float32:
		return float64(tv), true
	case float64:
		return tv, true
	case json.Number:
		val, err := tv.Float64()
		if err == nil {
			return val, true
		}
	}
	return 0, false
}

func MustBool(v any, defaultval ...bool) bool {
	val, ok := TryBool(v)
	if ok {
		return val
	}
	if len(defaultval) > 0 {
		return defaultval[0]
	}
	return false
}

func TryBool(v any) (ret bool, isbool bool) {
	if v == nil {
		return false, false
	}
	switch tv := v.(type) {
	case bool:
		return tv, true
	case string:
		//Attention:
		//   strconv.ParseBool() think "0","t","T","1","f","F" as bool,
		//   but it may not used as bool for outer func.
		//	 So only return value that must be used as bool.
		switch tv {
		case "true", "TRUE", "True":
			return true, true
		case "false", "FALSE", "False":
			return false, true
		default:
			return false, false
		}
	default:
		return false, false
	}
}
