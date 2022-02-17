package card

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		rank int32
		suit string
		want *Card
	}{
		{
			rank: Ranks()[0],
			suit: Suits()[0],
			want: &Card{Rank: Ranks()[0], Suit: Suits()[0], TextureKey: fmt.Sprintf("%02d%s", Ranks()[0], Suits()[0])},
		},
		{
			rank: Ranks()[0],
			suit: Empty,
			want: &Card{Rank: Ranks()[0], Suit: Empty, TextureKey: Empty},
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d/%s equals %v", tt.rank, tt.suit, tt.want)
		t.Run(testname, func(t *testing.T) {
			c := New(tt.rank, tt.suit)
			if !cmp.Equal(c, tt.want) {
				t.Errorf("got %v, want %v", c, tt.want)
			}
		})
	}
}

func TestCompareOverlappingSuit(t *testing.T) {
	var tests = []struct {
		src, dst *Card
		want     bool
	}{
		{
			src:  &Card{Rank: 1, Suit: "d"},
			dst:  &Card{Rank: 1, Suit: "d"},
			want: false,
		},
		{
			src:  &Card{Rank: 1, Suit: "d"},
			dst:  &Card{Rank: 1, Suit: "c"},
			want: true,
		},
		{
			src:  &Card{Rank: 1, Suit: "c"},
			dst:  &Card{Rank: 1, Suit: "s"},
			want: false,
		},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v/%v", tt.src, tt.dst)
		t.Run(testname, func(t *testing.T) {
			result := tt.src.CompareOverlappingSuit(tt.dst)
			if result != tt.want {
				t.Errorf("got %t, want %t", result, tt.want)
			}
		})
	}
}
