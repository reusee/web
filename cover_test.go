// +build cover

package web

/*

This test needs a modified testing package.

Add the following function to testing package to gain access to testing.Cover:

func GetCover() Cover {
  return cover
}

run with: go test -cover -tags cover

*/

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ret := m.Run()

	pt("------ not covered ------\n")
	cover := testing.GetCover()
	for filename, blocks := range cover.Blocks {
		for i, block := range blocks {
			count := cover.Counters[filename][i]
			if count > 0 {
				continue
			}
			pt("file %s line %d col %d\n", filename, block.Line0, block.Col0)
		}
	}

	os.Exit(ret)
}
