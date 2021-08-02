# MCsniperGO

## Usage

> This sniper is in it's *beta* stage, meaning bugs should be expected.

### Easy installation

- download the latest release for your operating system in the [releases](https://github.com/Kqzz/MCsniperGO/releases/) tab
- move that file to it's own folder (recommended, issues may arise if not done)
- run that file once, 2 files will appear in that folder, `accounts.txt` and `config.toml`. config.toml
  - on windows, you can double click the executable file to run.
  - on macos, you must [open a terminal in the folder you moved the binary to](https://www.stugon.com/open-terminal-in-current-folder-location-mac/) and run `./binary_name_here` in the terminal
  - on linux, you must open a terminal and use `cd` to navigate to the folder where the binary is located, you may have to run `chmod +x ./binary_name_here`, and then run `./binary_name_here`
- open accounts.txt and add your accounts according to the formatting specified [below](https://github.com/Kqzz/MCsniperGO#accounts-file-format)
- run the sniper again with the same commands used before
  - for windows, you can double click the executable
  - for macos, you will have to open the terminal in the correct directory (as shown above) and run `./binary_name_here`
  - for linux, you will have to open a terminal and use `cd` to navigate to the correct directory, then run `./binary_name_here`
- the sniper will now prompt you for a username and offset. enter those and then the sniper will authenticate (this is run 8 hours before snipe) and then count down.


### Compiling from source (not recommended)

install [git](https://git-scm.com/) and [go](https://golang.org/dl/) 1.16 or later.
```sh
git clone https://github.com/Kqzz/MCsniperGO
cd MCsniperGO
go run .

```

## Accounts file format

You have to put accounts in accounts.txt before sniping, so follow the format shown below to add your accounts

> in the code block shown below, replace any words fully capitalized with actual values, just leave the other words as is.

```txt
# You can comment out lines by including a # at the start

## MOJANG ACCOUNT

EMAIL:PASS
### OR (if you have security questions)
EMAIL:PASS:ANSWER:ANSWER:ANSWER

## Microsoft account
## IT IS NOT RECOMMENDED TO USE THIS, INSTEAD MANUALLY GRAB THE BEARER TOKEN AND USE THE BEARER METHOD.

EMAIL:PASS:ms
### OR FOR PRENAME
EMAIL:PASS:prename

## Manual bearer
BEARERHERE:bearer

# You can add :prename for a prename snipe w/ manual bearer

```

## MCsniperPY vs MCsniperGO

```txt

[MCsniperGO]

### PROS ###

* faster (go rather than py)
* microsoft auth support
* manual bearer
* prename sniping

### CONS ###

* low or no support
* brand new, not proven to be good.
* not tested well, very tiny userbase.
* basically in beta or alpha

[MCsniperPY]

### PROS ###

* proven to get names
* very well trusted and old
* huge discord server who can offer support

### CONS ###

* slower than MCsniperGO
* no MS auth
* no prename
* no manual bearer, for the few people who want that.

```
