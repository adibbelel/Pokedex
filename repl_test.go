package main

import (
  "testing"
)


func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
    {
		  input: " hello world ",
		  expected: []string{"hello", "world"},
	  },
	  {
		input: " xyz_vga 1, 2.3 4 ",
		expected: []string{"xyz_vga", "1,", "2.3", "4"},
	  },
  }

  for _, c := range cases {
    actual := cleanInput(c.input)
    
    if len(actual) != len(c.expected) {
      t.Errorf("oopsie, doesn't match the length: expected %d, got %d", len(c.expected), len(c.input))
      continue
    }

    for i := range actual {
      word := actual[i]
      expectedWord := c.expected[i]

      if word != expectedWord {
        t.Errorf("oopsie, words don't match")
      }
    }
  }
} 
