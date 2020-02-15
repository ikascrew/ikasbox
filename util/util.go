package util

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func SearchDirectory(d string, target []string) ([]string, error) {

	rtnFiles := make([]string, 0)
	files, err := ioutil.ReadDir(d)
	if err != nil {
		return nil, err
	}

	for _, elm := range files {
		if elm.IsDir() {
			childPath := filepath.Join(d, elm.Name())
			rtn, err := SearchDirectory(childPath, target)
			if err != nil {
				return nil, err
			}
			rtnFiles = append(rtnFiles, rtn...)
		} else {
			fName := elm.Name()
			appFlag := false
			if target != nil {
				for _, elm := range target {
					reg := regexp.MustCompile("(\\s" + elm + "$)")
					if reg.Match([]byte(fName)) {
						appFlag = true
						break
					}
				}
			} else {
				appFlag = true
			}
			if appFlag {
				name := filepath.Join(d, fName)
				rtnFiles = append(rtnFiles, name)
			}

		}
	}
	return rtnFiles, nil
}

func SortFiles(f []string) {
	sort.Slice(f, func(i, j int) bool {

		f1 := f[i]
		f2 := f[j]

		s1 := strings.Split(f1, string(filepath.Separator))
		s2 := strings.Split(f2, string(filepath.Separator))

		l1 := len(s1)
		l2 := len(s2)

		for idx := 0; idx < l1-1; idx++ {
			elm1 := s1[idx]
			if idx == l2-1 {
				return false
			}
			elm2 := s2[idx]
			if elm1 != elm2 {
				return elm1 < elm2
			}
		}

		if len(s1) == len(s2) {
			return f1 < f2
		}
		return l1 < l2
	})
}
