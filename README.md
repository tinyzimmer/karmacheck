# karmacheck
[![Build Status](https://www.travis-ci.com/tinyzimmer/karmacheck.svg?branch=master)](https://www.travis-ci.com/tinyzimmer/karmacheck)

Unofficial KarmaDecay checker for watching subreddits

Currently watches a subreddit (or list of subreddits) for new posts, and runs new submissions against KarmaDecay.
The program fetches the markdown comment from the KarmaDecay page and prints it to the terminal.

Pre-compiled binaries for Windows, Linux, and macOS can be found in the [releases](https://github.com/tinyzimmer/karmacheck/releases) section.

## Building

 - Tested on Go 1.10.2 windows/linux/macOS

#### Build

Using go:

```bash
$> go get github.com/tinyzimmer/karmacheck
```

From git:

```bash
$> git clone https://github.com/tinyzimmer/karmacheck
$> cd karmacheck
$> go build .
```

## Usage
```bash
$> ./karmacheck
Usage of karmacheck:
  -d    Debug
  -s string
        Comma separated list of subs to watch
```

#### Example
```powershell
PS C:\Users\tinyzimmer\Desktop\Development\karmacheck> .\karmacheck.exe -s funny
2018/07/03 08:22:27 Checking KarmaDecay for: r/funny/comments/8vt1il/praise_the_ol_mighty/
2018/07/03 08:22:28 KarmaDecay could not find any similar posts
```

```powershell
PS C:\Users\tinyzimmer\Desktop\Development\karmacheck> .\karmacheck.exe -s peoplefuckingdying
2018/07/03 08:23:13 Checking KarmaDecay for: r/PeopleFuckingDying/comments/8vsv4a/woman_is_consumed_alive_by_vicious_animals/
2018/07/03 08:23:17 Found matches. Below is the reddit comment text.
Anyone seeking more info might also check here:

title | points | age | /r/ | comnts
:--|:--|:--|:--|:--
[They have accepted me as their own](http://www.reddit.com/r/trashpandas/comments/5ko9v3/they_have_accepted_me_as_their_own/) ^**B** | 722 | 1^yr | trashpandas | 47
[I have no idea what would compel somebody to do this](http://www.reddit.com/r/WTF/comments/q5yuo/i_have_no_idea_what_would_compel_somebody_to_do/) ^**B** | 819 | 6^yrs | WTF | 220
[I honestly don't know...](http://www.reddit.com/r/WTF/comments/nn07n/i_honestly_dont_know/) ^**B** | 1133 | 6^yrs | WTF | 369
[This could end badly.](http://www.reddit.com/r/WTF/comments/t5dsg/this_could_end_badly/) | 40 | 6^yrs | WTF | 8
[Hey! Let's feed some raccoons.](http://www.reddit.com/r/WTF/comments/sfron/hey_lets_feed_some_raccoons/) | 176 | 6^yrs | WTF | 22
[Raccoon Log - Day 17 - Acceptance is finally imminent.](http://www.reddit.com/r/funny/comments/28qgh2/raccoon_log_day_17_acceptance_is_finally_imminent/) | 437 | 4^yrs | funny | 36
[She's got a thing for raccoons.](http://www.reddit.com/r/WTF/comments/1debwq/shes_got_a_thing_for_raccoons/) | 426 | 5^yrs | WTF | 55
[me irl](http://www.reddit.com/r/me_irl/comments/48meel/me_irl/) | 413 | 2^yrs | me_irl | 13
[Maybe trying meth on thanksgiving wasn't the best idea...](http://www.reddit.com/r/WTF/comments/2wxhx0/maybe_trying_meth_on_thanksgiving_wasnt_the_best/) | 733 | 3^yrs | WTF | 77

*[Source: karmadecay](http://karmadecay.com/r/PeopleFuckingDying/comments/8vsv4a/woman_is_consumed_alive_by_vicious_animals/) (B = bigger)*
```

## TODO

 - The examples are off with the most current version
 - Just need to write the bit that can post comments now, torn between rolling a quick one or going all out with [graw](https://github.com/turnage/graw)
 - Write deployment examples, probably ansible or terraform
