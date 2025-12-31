package main

import (
	"fmt"
	"testing"
)

func TestBadWordReplace(t *testing.T) {
	cases := []struct {
		body  string
		extra string
		want  string
	}{
		{
			body: "I had something interesting for breakfast",
			want: "I had something interesting for breakfast",
		},
		{
			body:  "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			extra: "this should be ignored",
			want:  "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			body: "I really need a kerfuffle to go to bed sooner, Fornax !",
			want: "I really need a **** to go to bed sooner, **** !",
		},
	}
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			cleaned := getCleanedBody(c.body, badWords)
			if cleaned != c.want {
				t.Errorf("Expected %s, got %s", c.want, cleaned)
			} else {
				t.Logf("Expected %s, got %s", c.want, cleaned)
			}
		})
	}
}
