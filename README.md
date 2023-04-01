# primeweb

A prime number bot, but as a [website][].

[website]: https://primes.today/

## How it Works

This was hacked together somewhat quickly as the Twitter bot looks imminent to
shut down. There's a few components:

- [`main.go`](./main.go): a small script which repurposes parts of [primebot][]
  to generate the prime numbers, as well as backfill and import features
- [Hugo](https://gohugo.io): used to build the primes static site
- [`bot.yml`](./.github/workflows/bot.yml): a workflow which builds on an hourly
  cadence to generate the next prime number in sequence
- [`gh-pages.yml`](./.github/workflows/gh-pages.yml): a workflow which builds
  the github pages deployment of the static site on any push to `main`

State is written to the [`_current`](./_current) file, and is used to determine
the last prime number that was generated.

[primebot]: https://github.com/primes-today/primebot

## License

[MIT](./LICENSE)
