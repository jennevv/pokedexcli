package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello   world    ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "helloworld",
			expected: []string{"helloworld"},
		},
		{
			input:    "  helloworld  ",
			expected: []string{"helloworld"},
		},
		{
			input:    "water    toes      ",
			expected: []string{"water", "toes"},
		},
		{
			input:    "    water    toes",
			expected: []string{"water", "toes"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		if len(actual) != len(c.expected) {
			t.Errorf("length actual slice (%v) != length of expected slice (%v)", len(actual), len(c.expected))
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("expected word %s does not match %s", expectedWord, word)
			}
		}
	}
}
