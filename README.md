# Auctioneer

<div align="center">
  <img width="300" src="https://assets.worldofwarcraft.com/static/components/Logo/Logo-wow-sitenav.596840db77b4d485a44d65e897e3de57.png" alt="Связной">
</div>

---

[![buddy pipeline](https://app.buddy.works/dog-sky/auctioneer/pipelines/pipeline/299664/badge.svg?token=84ccc34b677f07df458235c07be06800aea9d16ed9df2d1b390863e7e1e97c3d "buddy pipeline")](https://app.buddy.works/dog-sky/auctioneer/pipelines/pipeline/299664)

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
`/api/v1/auc_search?item_name=Гаррош&region=eu&realm_name=Страж%20Смерти` | The result of the search for the requested item


## Getting Started 👲

1. Clone project
2. Create a new branch from _origin/develop_
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
