package scraper

import (
	"os"
	"strconv"
	"testing"
)

func TestScraper(t *testing.T) {
	os.Chdir("../..")
	scr := New(nil)
	rssFeed, err := scr.ReadFeedFromFile(BlogXML)
	if err != nil {
		t.Errorf("Error reading feed %v", err)
	}
	tests := map[string]struct {
		want string
		got  string
	}{
		"channel title": {
			want: "Boot.dev Blog",
			got:  rssFeed.Channel.Title,
		},
		"channel description": {
			want: "Recent content on Boot.dev Blog",
			got:  rssFeed.Channel.Description,
		},

		"first title": {
			want: "You're Not Qualified to Have an Opinion on TDD",
			got:  rssFeed.Channel.Items[0].Title,
		},
		"last title": {
			want: "Secure Random Numbers in Node.js",
			got:  rssFeed.Channel.Items[len(rssFeed.Channel.Items)-1].Title,
		},
		"items length": {
			want: "349",
			got:  strconv.Itoa(len(rssFeed.Channel.Items)),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.want != tc.got {
				t.Errorf("want: %s, got: %s", tc.want, tc.got)
			}
		})
	}

}
