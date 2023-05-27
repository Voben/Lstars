package lstarPack

// check if a element contains a slice or not
func slice_contains(s []string, ele string) bool{
	for _, element := range s{
		if ele == element{
			return true
		}
	}
	return false
}

// check if a nested map contains a string
func nest_map_contains(s map[string]map[string]string, ele string) bool{
	_, val := s[ele]
	return val 
}

// check if a map contains a string
func map_contains(s map[string]string, ele string) bool{
	_, val := s[ele]
	return val 
}
