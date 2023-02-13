---
title: About
date: 2023-02-11T07:00:00-0800
---

_Every prime number, eventually. (Or the heat death of the universe; whichever happens first.)_

"primes" is a prime number listing bot; the next sequential prime number is
revealed every hour.

This website grew out of a [Twitter bot][twitter] which was launched in October
2013, where it lived exclusively until January of 2018 when a [Mastodon
bot][mastodon] was also launched. In February of 2023, changes to Twitter rules
made this longevity of the bot uncertain, prompting this website to be built.

The Twitter and Mastodon bots are stateless, except in-memory: when they start
up, the fetch the last-posted prime from their respective accounts, and post the
next prime at the appropriate time. In those cases the first post, `2`, was
manually sent via the web client, and then the bot was started.

When this site was launched, the [first prime][first] was backdated to match the
first tweet made on the bot's Twitter account, and the rest was populated from a
partial export of Twitter data, which covered all primes from 2020 October 27 at
21:59:06 GMT ([768,409][]), up through 2023 February 08 at 04:59:42 GMT
([1,042,897][]). This partial import was due to Twitter not supporting export of
full history. To fill those gaps, all other numbers were generated to smear
evenly across the missing gap between the [first][] and that export.

This is meant to keep the project alive as social networks come and go, and
their rules change.

[twitter]: https://twitter.com/_primes_
[mastodon]: https://botsin.space/@primes
[first]: {{< ref "primes/2" >}}
[768,409]: {{< ref  "primes/768409" >}}
[1,042,897]: {{< ref "primes/1042897" >}}
