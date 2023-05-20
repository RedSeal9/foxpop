package foxpop

import (
	"strconv"
	"strings"
)

func StringifyUserJS(data Data) string {
	var sb strings.Builder
	for _, e := range data.Entries {
		objt := typeof(e.Value)
		sep := ""
		if objt == "string" {
			sep = "\""
		}

		sb.WriteString("user_pref(\"" + e.Name + "\", " + sep + unfuck(e.Value) + sep + ");")

		sb.WriteRune('\n')
	}
	return sb.String()
}

func unfuck(v interface{}) string {
	switch v.(type) {
	case int:
		return strconv.Itoa(int(v.(int)))
	case int64:
		return strconv.Itoa(int(v.(int64)))
	case bool:
		return strconv.FormatBool(v.(bool))
	case string:
		return string(v.(string))
	}
	return ""
}

func typeof(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case bool:
		return "bool"
	case string:
		return "string"
	default:
		return ""
	}
}
