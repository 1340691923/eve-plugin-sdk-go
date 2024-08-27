package util

import "github.com/hashicorp/go-version"

func LessThan(v1, v2 string) bool {
	ver1, err1 := version.NewVersion(v1)
	if err1 != nil {
		return false
	}
	ver2, err2 := version.NewVersion(v2)
	if err2 != nil {
		return false
	}

	return ver1.LessThan(ver2)
}
 
