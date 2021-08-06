package fmath

import "testing"

func TestAdd(t *testing.T) {
	t.Log(Add(0.1, 0.05))
	t.Log(Add(0.1, -0.05))
}

func TestSub(t *testing.T) {
	t.Log(Sub(0.05, 0.05))
}

func TestMul(t *testing.T) {
	t.Log(Mul(0.05, 0.05))
}

func TestDiv(t *testing.T) {
	t.Log(Div(0.05, 0.035))
}

func TestRound(t *testing.T) {
	t.Log(Round(0.054, 2))
}
