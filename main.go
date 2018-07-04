/**
    This file is part of KarmaCheck.

    KarmaCheck is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    KarmaCheck is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with KarmaCheck.  If not, see <http://www.gnu.org/licenses/>.
**/

package main

import (
	"flag"
	"log"
	"os"
)

const (
	EXIT_NO_ARGS                         = 1
	EXIT_INVALID_SUBREDDIT               = 2
	EXIT_DEAD_TRACKER                    = 3
	REDDIT_URL                           = "https://www.reddit.com"
	KARMA_DECAY_URL                      = "https://www.karmadecay.com"
	RSS_ARG                              = ".rss"
	RSS_URL_FORMAT                       = "%s/r/%s/new/%s"
	REQUEST_AGENT_HEADER                 = "User-Agent"
	REQUEST_AGENT                        = "RepostBot by u/jews4beer"
	NO_CONTENT_ERROR                     = "KarmaDecay could not locate any media in the post"
	NO_SIMILAR_POSTS_ERROR               = "KarmaDecay could not find any similar posts"
	MALFORMED_URL_ERROR                  = "Malformed URL: %s"
	NO_SUBREDDIT_ERROR                   = "Invalid subreddit"
	KARMA_DECAY_NO_CONTENT_STRING        = "Unable to find an image"
	LOCAL_FOUND_MATCHES_MESSAGE          = "Found matches. Below is the reddit comment text."
	LOCAL_BELOW_CONFIDENCE_MESSAGE       = "KarmaDecay response scored below the confidence threshold"
	MARKDOWN_SEARCH_REGEX                = "Anyone[^<]*"
	MARKDOWN_LINK_REGEX                  = "\\[.*\\]\\(.*\\)"
	MARKDOWN_VALID_CHECK                 = "[Source: karmadecay]"
	FEEDTRACKER_CHECKED_ENTRIES_MAX      = 100
	REDDIT_CHECK_SLEEP_TIME              = 10
	KARMA_DECAY_SLEEP_TIME               = 3
	CHECK_TRACKER_SLEEP_TIME             = 3
	KARMA_DECAY_COMMENT_LINKS_CONFIDENCE = 2
)

var (
	subreddits = flag.String("s", "", "Comma separated list of subs to watch")
)

func main() {

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(EXIT_NO_ARGS)
	}
	flag.Parse()
	subs, err := checkSubs(subreddits) // make sure subreddit is defined
	if err != nil {
		log.Println(err)
		flag.Usage()
		os.Exit(EXIT_INVALID_SUBREDDIT)
	}

	activeTrackers := []FeedTracker{}
	for _, sub := range subs {
		// Create a Reddit feed tracker and run it
		tracker := NewTracker(sub)
		activeTrackers = append(activeTrackers, tracker)
		go tracker.Run()
	}

	ManageTrackers(activeTrackers)

}
