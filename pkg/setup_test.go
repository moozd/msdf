package msdf

import (
	"io"
	"log"
	"os"
	"testing"
)

func TestMain(t *testing.M) {
	log.SetOutput(io.Discard)
	code := t.Run()
	os.Exit(code)

}
