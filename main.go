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
	// Exit Codes
	EXIT_NO_ARGS           = 1
	EXIT_INVALID_SUBREDDIT = 2
	EXIT_DEAD_TRACKER      = 3

	// URL stuff
	REDDIT_URL      = "https://www.reddit.com"
	KARMA_DECAY_URL = "https://www.karmadecay.com"

	// Request Agent
	REQUEST_AGENT_HEADER = "User-Agent"
	REQUEST_AGENT        = "KarmaCheck/v0.7 by u/jews4beer"

	// Error strings
	NO_CONTENT_ERROR       = "KarmaDecay could not locate any media in the post"
	NO_SIMILAR_POSTS_ERROR = "KarmaDecay could not find any similar posts"
	EMPTY_SUBREDDIT_ERROR  = "The subreddit does not appear to have any posts"
	MALFORMED_URL_ERROR    = "Malformed URL: %s"
	NO_SUBREDDIT_ERROR     = "Invalid subreddit"
	FAILED_INIT_ERROR      = "ERROR: Failed to initiate tracker for sub r/%s: %s\n"
	FAILED_POLL_ERROR      = "ERROR: Error polling results from sub r/%s: %s"
	DEAD_TRACKER_ERROR     = "One or more trackers have stopped running. Exiting."

	// KarmaDecay scrape strings
	KARMA_DECAY_NO_CONTENT_STRING = "Unable to find an image"

	// Local print messages
	LOCAL_FOUND_MATCHES_MESSAGE    = "Found matches. Below is the reddit comment text."
	LOCAL_BELOW_CONFIDENCE_MESSAGE = "KarmaDecay response scored below the confidence threshold"

	// KarmaDecay regex search strings
	MARKDOWN_SEARCH_REGEX = "Anyone[^<]*"
	MARKDOWN_LINK_REGEX   = "\\[.*\\]\\(.*\\)"
	MARKDOWN_VALID_CHECK  = "[Source: karmadecay]"

	// Other defaults
	KARMA_DECAY_COMMENT_LINKS_CONFIDENCE = 2
)

var (
	dryRun     = flag.Bool("d", false, "Dry Run. Do not reply to posts.")
	subreddits = flag.String("s", "", "Comma separated list of subs to watch")
	config     = flag.String("c", "./bot.agent", "Path to Bot Configuration.")
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

	StartRedditSession(*config, subs, *dryRun)

}
