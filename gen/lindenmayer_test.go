package gen

import (
	"testing"
)

func TestRewrite(t *testing.T) {
	for i, tt := range []struct {
		axiom       string
		productions map[rune]string
		want        string
	}{
		{
			axiom:       "",
			productions: map[rune]string{},
			want:        "",
		},
		{
			axiom: "F-F-F-F",
			productions: map[rune]string{
				'F': "F-F+F+FF-F-F+F",
			},
			want: "F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F-F-F+F+FF-F-F+F",
		},
	} {
		lsystem := Lsystem{
			Axiom:       tt.axiom,
			Productions: tt.productions,
		}
		got := lsystem.rewrite(lsystem.Axiom)
		if got != tt.want {
			t.Errorf("%d): got %s want %s", i, got, tt.want)
		}
	}
}
