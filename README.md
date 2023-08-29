<h1 align="center">
  <img src="https://i.imgur.com/c9OxLeM.png" alt="MCsniperGO"></img>
</h1>

> By Kqzz ~ [Discord](https://discord.gg/mcsnipergo-734794891258757160)

## Usage

- [Install go](https://go.dev/dl/)
- Download or clone MCsniperGO repository 
- open MCsniperGO folder in your terminal / cmd
- put your prename accounts (no claimed username) in [`gc.txt`](#accounts-formatting) and your normal accounts in [`ms.txt`](#accounts-formatting)
- put proxies into `proxies.txt` in the format `user:pass@ip:port` (there should NOT be 4 `:` in it as many proxy providers provide it as)
- run `go run ./cmd/cli`
- enter username + [claim range](#claim-range)
- wait, and hope you claim the username!

## Claim Range
Use the following Javascript bookmarklet in your browser to obtain the droptime while on `namemc.com/search?q=<username>`:

```js
javascript:(function(){function parseIsoDatetime(dtstr) {
    return new Date(dtstr);
};

startElement = document.getElementById('availability-time');
endElement = document.getElementById('availability-time2');

start = parseIsoDatetime(startElement.getAttribute('datetime'));
end = parseIsoDatetime(endElement.getAttribute('datetime'));

para = document.createElement("p");
para.innerText = Math.floor(start.getTime() / 1000) + '-' + Math.ceil(end.getTime() / 1000);

endElement.parentElement.appendChild(para);})();

```

If 3name.xyz has a lower length claim range for a username I would recommend using that, you can get the unix droptime range with this bookmarklet on `3name.xyz/name/<name>`

```js
javascript: (function() {
    startElement = document.getElementById('lower-bound-update');
    endElement = document.getElementById('upper-bound-update');
    para = document.createElement("p");
    para.innerText = startElement.getAttribute('data-lower-bound') + '-' + endElement.getAttribute('data-upper-bound');
    endElement.parentElement.appendChild(para)
})()
```

## accounts formatting

`gc.txt` and `ms.txt`
```txt
EMAIL:PASSWORD
```

## understanding logs

Each request made to change your username will return a 3 digit HTTP status code, the meanings are as follows:

- 400 / 403: Failed to claim username (will continue trying)
- 401: Unauthorized (restart claimer if it appears)
- 429: Too many requests (add more proxies if this occurs frequently)
