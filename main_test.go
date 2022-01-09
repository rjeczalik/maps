package objects_test

import (
	"flag"
	"os"
	"testing"
)

var updateGolden = flag.Bool("update-golden", false, "Updates golden files")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}
