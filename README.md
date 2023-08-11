<h1 align="center">
  <img src="https://i.imgur.com/c9OxLeM.png" alt="MCsniperGO"></img>
</h1>

> By Kqzz ~ [Discord](https://discord.gg/mcsnipergo-734794891258757160)

## Usage

- [Install go](https://go.dev/dl/)
- open MCsniperGO folder in terminal / cmd
- put your prename accounts in [`gc.txt`](#accounts-formatting) (proxy support + normal MS support soon)
- run `go run ./cmd/cli`
- Enter username + [claim range](#claim-range)

## Claim Range
Use the following Javascript bookmarklet in your browser to obtain the droptime while on `namemc.com/search?q=<username>`:

```js
javascript:(function()%7Bfunction%20parseIsoDatetime(dtstr)%20%7B%0A%20%20%20%20return%20new%20Date(dtstr)%3B%0A%7D%3B%0A%0AstartElement%20%3D%20document.getElementById('availability-time')%3B%0AendElement%20%3D%20document.getElementById('availability-time2')%3B%0A%0Astart%20%3D%20parseIsoDatetime(startElement.getAttribute('datetime'))%3B%0Aend%20%3D%20parseIsoDatetime(endElement.getAttribute('datetime'))%3B%0A%0Apara%20%3D%20document.createElement(%22p%22)%3B%0Apara.innerText%20%3D%20Math.floor(start.getTime()%20%2F%201000)%20%2B%20'-'%20%2B%20Math.ceil(end.getTime()%20%2F%201000)%3B%0A%0AendElement.parentElement.appendChild(para)%3B%7D)()%3B
```

## accounts formatting

```txt
EMAIL:PASSWORD
```

## understanding logs

Each request made to change your username will return a 3 digit HTTP status code, the meanings are as follows:

- 400 / 403: Failed to claim username (will continue trying)
- 401: Unauthorized (restart claimer if it appears)
- 429: Too many requests (add more proxies if this occurs frequently)
