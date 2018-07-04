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
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

// Type definitions for Reddit RSS feed parsing

type Feed struct {
	Namespace string   `xml:"xmlns,attr"`
	Category  Category `xml:"category"`
	Updated   string   `xml:"updated"`
	Icon      string   `xml:"icon"`
	Id        string   `xml:"id"`
	Links     []Link   `xml:"link"`
	Title     string   `xml:"title"`
	Entries   []Entry  `xml:"entry"`
}

type Category struct {
	Term  string `xml:"term,attr"`
	Label string `xml:"label,attr"`
}

type Link struct {
	Rel  string `xml:"rel,attr"`
	Href string `xml:"href,attr"`
	Type string `xml:"type,attr"`
}

type Entry struct {
	Author  Author `xml:"author"`
	Content string `xml:"content"`
	Id      string `xml:"id"`
	Link    Link   `xml:"link"`
	Updated string `xml:"updated"`
	Title   string `xml:"title"`
}

type Author struct {
	Name string `xml:"name"`
	Uri  string `xml:"uri"`
}

type FeedTracker struct {
	Subreddit      string
	Running        bool
	Debug          bool
	CheckedEntries []Entry
}

func ManageTrackers(trackers []*FeedTracker) {
	var active int
	var tracker *FeedTracker
	for {
		active = 0
		for _, tracker = range trackers {
			if tracker.Debug {
				log.Printf("DEBUG: Checking tracker: r/%s\n", tracker.Subreddit)
			}
			if tracker.Running {
				active += 1
			}
		}
		if active != len(trackers) {
			log.Println("One or more trackers have stopped running. Exiting.")
			os.Exit(EXIT_DEAD_TRACKER)
		}
		time.Sleep(time.Second * CHECK_TRACKER_SLEEP_TIME)
	}
}

func NewTracker(sub string, debug bool) (f FeedTracker) {
	f.Subreddit = sub
	f.Running = true
	f.Debug = debug
	return
}

func (f *FeedTracker) Run() {

	if f.Debug {
		log.Println("DEBUG: Whitelisting pre-existing entries")
	}
	err := f.InitializeCheckedEntries() // skip over pre-existing posts
	if err != nil {
		f.Die()
		return
	}
	log.Printf("Subreddit tracker started for r/%s\n", f.Subreddit)

	for {
		if f.Debug {
			log.Printf("DEBUG: Polling subreddit: r/%s\n", f.Subreddit)
		}
		data, err := f.GetLatestSubmissions() // get latest reddit posts
		if err != nil {
			log.Printf("ERROR: Error polling results from sub r/%s: %s", f.Subreddit, err.Error())
			f.Die()
			return
		}
		// iterate through the response and check each entry

		for _, entry := range data {
			if !f.IsRecorded(entry) {
				res, err := checkKarmaDecay(entry)
				if err != nil {
					log.Println(err)
				} else {
					if kdIsConfident([]byte(res)) {
						log.Println(LOCAL_FOUND_MATCHES_MESSAGE)
						fmt.Println(res)
					} else {
						log.Println(LOCAL_BELOW_CONFIDENCE_MESSAGE)
					}
				}
				f.RecordEntry(entry)
				time.Sleep(time.Second * KARMA_DECAY_SLEEP_TIME)
			}
		}

		time.Sleep(time.Second * REDDIT_CHECK_SLEEP_TIME)

	}
}

func (f *FeedTracker) InitializeCheckedEntries() (err error) {
	data, err := f.GetLatestSubmissions()
	if err != nil {
		log.Printf("ERROR: Failed to initiate tracker for sub r/%s: %s\n", f.Subreddit, err.Error())
		return
	}
	for _, entry := range data {
		f.RecordEntry(entry)
	}
	return
}

func (f *FeedTracker) RecordEntry(entry Entry) {
	if len(f.CheckedEntries) >= FEEDTRACKER_CHECKED_ENTRIES_MAX {
		f.RemoveCheckedEntry(0)
	}
	f.AddCheckedEntry(entry)
	if f.Debug {
		log.Printf("DEBUG: Recorded Entry: %s", entry.Title)
	}
}

func (f *FeedTracker) GetLatestSubmissions() (entries []Entry, err error) {

	// Poll the reddit RSS feed for the latest submissions from a subreddit

	url := fmt.Sprintf(RSS_URL_FORMAT, REDDIT_URL, f.Subreddit, RSS_ARG)
	resp, err := getUrl(url)
	if err != nil {
		f.Die()
		return
	}
	feed := Feed{}
	xml.Unmarshal(resp, &feed)
	if len(feed.Entries) == 0 {
		f.Die()
		err = errors.New(EMPTY_SUBREDDIT_ERROR)
		return
	}
	entries = feed.Entries[:10]
	return

}

func (f *FeedTracker) AddCheckedEntry(entry Entry) {
	appended := append(f.CheckedEntries, entry)
	f.CheckedEntries = appended
}

func (f *FeedTracker) RemoveCheckedEntry(i int) {
	f.CheckedEntries = append(f.CheckedEntries[:i], f.CheckedEntries[i+1:]...)
}

func (f *FeedTracker) IsRecorded(entry Entry) bool {
	for _, item := range f.CheckedEntries {
		if item.Link.Href == entry.Link.Href {
			return true
		}
	}
	return false
}

func (f *FeedTracker) Die() {
	f.Running = false
}
