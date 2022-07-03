package stargazer

import "os"

func isFile(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.Mode().IsRegular()
}

func isDir(path string) bool {
	st, err := os.Stat(path)
	if err != nil {
		return false
	}
	return st.IsDir()
}
