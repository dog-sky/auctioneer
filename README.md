# Auctioneer

<div align="center">
  <img width="300" src="https://assets.worldofwarcraft.com/static/components/Logo/Logo-wow-sitenav.596840db77b4d485a44d65e897e3de57.png" alt="Связной">
</div>

---

[![buddy pipeline](https://app.buddy.works/dog-sky/auctioneer/pipelines/pipeline/299664/badge.svg?token=84ccc34b677f07df458235c07be06800aea9d16ed9df2d1b390863e7e1e97c3d "buddy pipeline")](https://app.buddy.works/dog-sky/auctioneer/pipelines/pipeline/299664)
![Test and Lint](https://github.com/dog-sky/auctioneer/workflows/Go/badge.svg?branch=main)
[![codecov](https://codecov.io/gh/dog-sky/auctioneer/branch/main/graph/badge.svg?token=SADKGY8ORK)](https://codecov.io/gh/dog-sky/auctioneer)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=dog-sky_auctioneer&metric=alert_status)](https://sonarcloud.io/dashboard?id=dog-sky_auctioneer)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/dog-sky/auctioneer)
[![TelegramBot](https://img.shields.io/badge/t.me-TelegramBot-blue)](https://t.me/AuctionHouseWoWBot)

---

<div align="center">
  <a href="#api-methods-">API Methods</a> •
  <a href="#getting-started-">Getting Started</a> •
  <a href="#local-development-">Local Development</a> •
  <a href="#helpful-links-">Helpful Links</a>
</div>

## API Methods 🚀

Route | Expected result
------------ | -------------
`/api/v1/auc_search?item_name={{string}}&region={{string}}&realm_name={{string}}` | The result of the search for the requested item
`/api/v1/item_media/{{int}}` | Awailable media for requested item ID


`auc search` query params | Description
------------ | -------------
`item_name` *required | Name of item to search
`region` *required | Region of realm. Choose `eu` or `us`
`realm_name` *required | Name of realm in the region


## Getting Started 👲

1. Clone project
2. Create a new branch from _origin/main_
3. Do `go mod download` in the root directory
4. Make the magic
5. Push branch and create merge request
6. Be happy with the results and do it again from the beginning

## Local Development 🚧
Export environment variables for example:
  ```
AUCTIONEER_BLIZZARD_EU_API_URL = "https://eu.api.blizzard.com"
AUCTIONEER_BLIZZARD_US_API_URL = "https://us.api.blizzard.com"
AUCTIONEER_BLIZZARD_AUTH_URL = "https://us.battle.net/oauth/token"
AUCTIONEER_BLIZZARD_CLIENT_SECRET = "your_client_secret"
AUCTIONEER_BLIZZARD_CLIENT_ID = "your_client_id"
AUCTIONEER_APP_PORT = ":8000"
AUCTIONEER_LOG_LEVEL = "DEBUG"
  ```

## Helpful Links 🤔

[World of Warcraft API](https://develop.battle.net/documentation/world-of-warcraft)
