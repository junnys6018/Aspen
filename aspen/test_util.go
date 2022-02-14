package main

import "testing"

type TestCase interface {
	Run(t *testing.T)
}
