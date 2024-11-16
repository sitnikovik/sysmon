package strings

// ToInterfaces maps command string arguments to interfaces for OS.
func ToInterfaces(ss []string) []interface{} {
	res := make([]interface{}, len(ss))
	for i, s := range ss {
		res[i] = interface{}(s)
	}

	return res
}
