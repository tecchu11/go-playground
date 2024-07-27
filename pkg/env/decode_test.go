package env_test

import (
	"go-playground/pkg/env"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	type S struct {
		Name string `env:"Name"`
		Age  int    `env:"Age"`
	}
	os.Setenv("Name", "tecchu11")
	os.Setenv("Age", "31")
	var s S
	env.Decode(&s)
	if s.Name != "tecchu11" {
		t.Fatalf("given name is %s", s.Name)
	}
	if s.Age != 31 {
		t.Fatalf("given age is %d", s.Age)
	}
}
