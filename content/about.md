---
title: About
date: 2023-02-11T07:00:00-0800
lastmod: 2023-02-25T21:23:00-0800
---

_Every prime number, eventually. (Or the heat death of the universe; whichever happens first.)_

"primes" is a prime number listing bot; the next sequential prime number is
revealed hourly.

This website grew out of a [Twitter bot][twitter] which was launched in October
2013, where it lived exclusively until January of 2018 when a [Mastodon
bot][mastodon] was also launched. In February of 2023, changes to Twitter rules
made this longevity of the bot uncertain, prompting this website to be built.

This is now considered to be the canonical version of the bot by its creator.

# How the social bots work 

The Twitter and Mastodon bots are stateless, except in-memory: when they start
up, the fetch the last-posted prime from their respective accounts, and post the
next prime at the appropriate time. In those cases the first post, `2`, was
manually sent via the web client, and then the bot was started.

They do not post on a known schedule, they instead check the time of the last
sent message, and post so long as it's been an hour since the last post. In
practice, this means that the posts are exactly an hour apart, but in cases
where a site or bot outage occurred, it can skew the next post and all
subsequent posts by the length of the outage.

# How this site was seeded

When this site was launched, the [first prime][first] was backdated to match the
first tweet made on the bot's Twitter account, and the rest was populated from a
partial export of Twitter data, which covered all primes from _2020 October 27
at 21:59:06 GMT_ ([768,409][]), up through _2023 February 08 at 04:59:42 GMT_
([1,042,897][]). This partial import was due to Twitter not supporting export of
full history. To fill those gaps, all other numbers were generated to smear
evenly across the missing gap between the [first][] and that export.

# Differences between the site and social bots

Unlike the social bots, this site is set to update every hour; however this
schedule is a request only: the actual may start some time after that request is
made. Ideally this means the bot posts hourly, but if the request to run an
update takes time to be fulfilled, it may post later.

This update is performed with [GitHub Actions][]; you can see the source that
builds this site [on GitHub][source].


[twitter]: https://twitter.com/_primes_
[mastodon]: https://botsin.space/@primes
[first]: {{< ref "primes/2" >}}
[768,409]: {{< ref  "primes/768409" >}}
[1,042,897]: {{< ref "primes/1042897" >}}
[GitHub Actions]: https://docs.github.com/en/actions
[source]: {{% param "repository.url" %}}
