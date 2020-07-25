package gen

import (
	"testing"
)

func TestRewriteN(t *testing.T) {
	for i, tt := range []struct {
		axiom       string
		productions map[rune]string
		n           int
		want        string
	}{
		{
			axiom:       "",
			productions: map[rune]string{},
			n:           0,
			want:        "",
		},
		{
			axiom: "F-F-F-F",
			productions: map[rune]string{
				'F': "F-F+F+FF-F-F+F",
			},
			n:    1,
			want: "F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F",
		},
		{
			axiom: "F-F-F-F",
			productions: map[rune]string{
				'F': "F-F+F+FF-F-F+F",
			},
			n:    2,
			want: "F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F+F-F+F+FF-F-F+FF-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F+F-F+F+FF-F-F+FF-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F+F-F+F+FF-F-F+FF-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F+F-F+F+FF-F-F+FF-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F+F-F+F+FF-F-F+F",
		},
	} {
		lsystem := Lsystem{
			Axiom:       tt.axiom,
			Productions: tt.productions,
		}
		got := lsystem.rewriteN(tt.axiom, tt.n)
		if got != tt.want {
			t.Errorf("%d): got %s want %s", i, got, tt.want)
		}
	}
}
