# vaultstats

This is a prototype/proof of concept to collect stat information on
github.com/hashicorp/vault issues

### Install:

1) clone this repo
2) run `make`

### Prerequisite:

Personal access token from here https://github.com/settings/tokens I think it just needs `public_repo`, `read:org` and `notifications`


    $ export GITHUB_API_TOKEN=""

### Usage:

    $ vaultstats stats


`stats` is the only command right now
