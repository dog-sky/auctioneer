# Auctioneer

<div align="center">
  <img width="300" src="https://assets.worldofwarcraft.com/static/components/Logo/Logo-wow-sitenav.596840db77b4d485a44d65e897e3de57.png" alt="Связной">
</div>

---

<div align="center">
  <a href="#api-methods-">API Methods</a> •
  <a href="#getting-started-">Getting Started</a> •
  <a href="#local-development-">Local Development</a> •
  <a href="#helpful-links-">Helpful Links</a>
</div>

## API Methods 🚀

<table>
  <thead>
    <tr>
      <th>Route</th>
      <th>Expected result</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>/api/v1/auc_search?item_name=Гаррош&region=eu&realm_name=Страж%20Смерти</code></td>
      <td>The result of the search for the requested item</td>
    </tr>
  </tbody>
</table>

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
AUCTIONEER_BLIZZARD_API_URL = "https://%s.api.blizzard.com"
AUCTIONEER_BLIZZARD_AUTH_URL = "https://us.battle.net/oauth/token"
AUCTIONEER_BLIZZARD_CLIENT_SECRET = "your_client_secret"
AUCTIONEER_BLIZZARD_CLIENT_ID = "your_client_id"
AUCTIONEER_APP_PORT = ":8000"
AUCTIONEER_LOG_LEVEL = "DEBUG"
  ```

## Helpful Links 🤔

[World of Warcraft API](https://develop.battle.net/documentation/world-of-warcraft)
