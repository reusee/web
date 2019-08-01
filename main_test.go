package web

import (
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ret := m.Run()
	time.Sleep(time.Second)
	os.Exit(ret)
}
