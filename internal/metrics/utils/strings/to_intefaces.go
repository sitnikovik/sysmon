package strings

// ToInterfaces maps command string arguments to interfaces for os
func ToInterfaces(ss []string) []interface{} {
	res := make([]interface{}, len(ss))
	for i, s := range ss {
		res[i] = interface{}(s)
	}

	return res
}
