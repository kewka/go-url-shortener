package test

import (
	"log"
	"os"
	"testing"
)

var env *Env

func TestMain(m *testing.M) {
	var err error
	env, err = NewEnv()
	if err != nil {
		log.Fatal(err)
	}
	code := m.Run()
	env.Cleanup()
	os.Exit(code)
}
