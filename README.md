# Nexus Hub
[![Build Status](https://github.com/mavryk-network/nexushub/workflows/build/badge.svg)](https://github.com/mavryk-network/nexushub/actions?query=branch%3Amaster+workflow%3A%22build%22)
[![made_with golang](https://img.shields.io/badge/made_with-golang-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbaking-bad%2Fnexushub.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbaking-bad%2Fnexushub?ref=badge_shield)

## Quickstart

### Sandbox

The simplest way is just to copy the `docker-compose.mavbox.yml` to your project.

Make sure you have the latest images and run the compose:
```
make sandbox-pull
make mavbox-sandbox
```
Sandbox UI is now available at http://localhost:8000


In order to stop or reset sandbox:
```
make sandbox-down
make sandbox-clear
```

## Read more

* [Configuration](./docs/configuration.md)
* [Developer docs](./docs/developer.md)


## Contact us
* Telegram: https://t.me/baking_bad_chat
* Twitter: https://twitter.com/TezosBakingBad
* Slack: https://tezos-dev.slack.com/archives/CV5NX7F2L
* Discord: https://discord.gg/RcPGSdcVSx


## About
This project is a fork of [BCDHUB](https://github.com/baking-bad/bcdhub) by baking-bad.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbaking-bad%2Fnexushub.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbaking-bad%2Fnexushub?ref=badge_large)