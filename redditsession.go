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
	"fmt"
	"log"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type RepostBot struct {
	bot    reddit.Bot
	dryRun bool
}

func (r *RepostBot) Post(p *reddit.Post) (err error) {
	res, e := checkKarmaDecay(*p)
	if e != nil {
		log.Println(e)
	} else {
		if kdIsConfident([]byte(res)) {
			log.Println(LOCAL_FOUND_MATCHES_MESSAGE)
			fmt.Println(res)
			if !r.dryRun {
				return r.bot.Reply(p.Name, res)
			} else {
				log.Println("Running dry-run mode. Skipping reply.")
			}
		} else {
			log.Println(LOCAL_BELOW_CONFIDENCE_MESSAGE)
		}
	}
	return
}

func StartRedditSession(subs []string, dryRun bool) {
	if bot, err := reddit.NewBotFromAgentFile("bot.agent", 0); err != nil {
		log.Fatal("Failed to create bot handle: ", err)
	} else {
		cfg := graw.Config{Subreddits: subs}
		handler := &RepostBot{
			bot:    bot,
			dryRun: dryRun,
		}
		log.Printf("Starting repost trackers for subs: %v\n", subs)
		if dryRun {
			log.Printf("Running in dry-run mode. Will not reply to posts.")
		}
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			fmt.Println("Failed to start graw run: ", err)
		} else {
			fmt.Println("graw run failed: ", wait())
		}
	}
}
