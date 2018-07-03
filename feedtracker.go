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
	"fmt"
	"log"
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
	Subreddit      *string
	CheckedEntries []Entry
}

func (f *FeedTracker) Run() {

	log.Println("Whitelisting pre-existing entries")
	f.InitializeCheckedEntries() // skip over pre-existing posts
	log.Println("Subreddit tracker started")

	for {

		log.Printf("Polling subreddit: r/%s\n", *f.Subreddit)
		data, err := f.GetLatestSubmissions() // get latest reddit posts
		if err != nil {
			log.Fatal(err)
		}
		// iterate through the response and check each entry

		for _, entry := range data {
			if !f.IsRecorded(entry) {
				res, err := checkKarmaDecay(entry)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(LOCAL_FOUND_MATCHES)
					fmt.Println(res)
				}
				f.RecordEntry(entry)
				time.Sleep(time.Second * KARMA_DECAY_SLEEP_TIME)
			}
		}

		time.Sleep(time.Second * REDDIT_CHECK_SLEEP_TIME)

	}
}

func (f *FeedTracker) InitializeCheckedEntries() {
	data, err := f.GetLatestSubmissions()
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range data {
		f.RecordEntry(entry)
	}
}

func (f *FeedTracker) RecordEntry(entry Entry) {
	if len(f.CheckedEntries) >= FEEDTRACKER_CHECKED_ENTRIES_MAX {
		f.RemoveCheckedEntry(0)
	}
	f.AddCheckedEntry(entry)
	log.Printf("Recorded Entry: %s", entry.Title)
}

func (f *FeedTracker) GetLatestSubmissions() (entries []Entry, err error) {

	// Poll the reddit RSS feed for the latest submissions from a subreddit

	url := fmt.Sprintf(RSS_URL_FORMAT, REDDIT_URL, *f.Subreddit, RSS_ARG)
	resp, err := getUrl(url)
	if err != nil {
		return
	}
	feed := Feed{}
	xml.Unmarshal(resp, &feed)
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
		if item.Title == entry.Title {
			return true
		}
	}
	return false
}
