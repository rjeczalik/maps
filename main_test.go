package objects_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var updateGolden = flag.Bool("update-golden", false, "Updates golden files")

func TestMain(m *testing.M) {
	flag.Parse()

	os.Exit(m.Run())
}

func Equal(rhs, lhs any) error {
	var v1, v2 any

	if err := reencode(rhs, &v1); err != nil {
		return err
	}

	if err := reencode(lhs, &v2); err != nil {
		return err
	}

	if !cmp.Equal(v1, v2) {
		return fmt.Errorf("rhs != lhs:\n%s", cmp.Diff(v1, v2))
	}

	return nil
}

func reencode(in, out any) error {
	p, err := json.Marshal(in)
	if err != nil {
		return err
	}

	return json.Unmarshal(p, out)
}
