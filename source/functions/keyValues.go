package functions

import (
	"Nosviak4/source"
	"fmt"
	"strconv"
	"strings"
)

// HandleKeyValues will convert the data into it's represented fields
func HandleKeyValues(args []string, method *source.Method) (map[string]interface{}, error) {
	kvs := make(map[string]interface{})

	// parses what has been provided
	for _, arg := range args {
		if !strings.HasPrefix(arg, source.MethodConfig.Attacks.KvPrefix) {
			continue
		}

		token := arg[len(source.MethodConfig.Attacks.KvPrefix):]

		kv, ok := method.Options.KeyValues[strings.Split(token, "=")[0]]
		if !ok || kv == nil || !strings.Contains(token, "=") {
			return make(map[string]interface{}), fmt.Errorf("key value error near %s", arg)
		}

		value := strings.Split(token, "=")[1]
		if len(value) > kv.StringMax {
			return make(map[string]interface{}), fmt.Errorf("key value error near %s", arg)
		}

		switch kv.Type {

		case "INT", "NUMBER":
			conv, err := strconv.Atoi(strings.Split(token, "=")[1])
			if err != nil || conv > kv.IntMax {
				return make(map[string]interface{}), fmt.Errorf("key value error near %s", arg)
			}

			kvs[strings.Split(token, "=")[0]] = conv

		case "BOOL":
			conv, err := strconv.ParseBool(strings.Split(token, "=")[1])
			if err != nil {
				return make(map[string]interface{}), fmt.Errorf("key value error near %s", arg)
			}

			kvs[strings.Split(token, "=")[0]] = conv

		case "STRING":
			kvs[strings.Split(token, "=")[0]] = strings.Split(token, "=")[1]
		}
	}

	// compares with our original config
	if len(kvs) == len(method.Options.KeyValues) {
		return kvs, nil
	}

	for key, value := range method.Options.KeyValues {
		if _, ok := kvs[key]; ok {
			continue
		}

		if value.Required {
			return make(map[string]interface{}), fmt.Errorf("missing required key value %s", key)
		}

		kvs[key] = value.Default
	}

	return kvs, nil
}