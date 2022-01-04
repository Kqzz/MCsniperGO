# MCsniperGO
This project was made possible by my [donators](https://kqzz.me/donators)

## Usage

> This sniper is in it's *beta* stage, meaning bugs should be expected.

### Easy installation

- download the latest release for your operating system in the [releases](https://github.com/Kqzz/MCsniperGO/releases/) tab
- move that file to it's own folder (recommended, issues may arise if not done)
- run that file once, 2 files will appear in that folder, `accounts.txt` and `config.toml`. `config.toml` can be ignored for now.
  - on windows, you can double click the executable file to run.
  - on macos, you must [open a terminal in the folder you moved the binary to](https://www.stugon.com/open-terminal-in-current-folder-location-mac/), run `chmod +x ./binary_name_here`, and then run`./binary_name_here` in the terminal
  - on linux, you must open a terminal and use `cd` to navigate to the folder where the binary is located, you may have to run `chmod +x ./binary_name_here`, and then run `./binary_name_here`
- open accounts.txt and add your accounts according to the formatting specified [below](https://github.com/Kqzz/MCsniperGO#accounts-file-format)
- run the sniper again with the same commands used before
  - for windows, you can double click the executable
  - for macos, you will have to open the terminal in the correct directory (as shown above) and run `./binary_name_here`
  - for linux, you will have to open a terminal and use `cd` to navigate to the correct directory, then run `./binary_name_here`
- the sniper will now prompt you for a username and offset. enter those and then the sniper will authenticate (this is run 8 hours before snipe) and then count down.

> A video guide will be made once this sniper is stabler

### Compiling from source (not recommended)

install [git](https://git-scm.com/) and [go](https://golang.org/dl/) 1.16 or later.
```sh
git clone https://github.com/Kqzz/MCsniperGO
cd MCsniperGO
go run .

```

## Accounts file format

You have to put accounts in accounts.txt before sniping, so follow the format shown below to add your accounts

> **in the code block shown below, replace any words fully capitalized with actual values, just leave the other words as is.**

```txt
# You can comment out lines by including a # at the start
# !! IMPORTANT !! ONLY REPLACE THE CAPITALIZED WORDS WITH THE CORRECT VALUE, EVERYTHING ELSE IS A "FLAG" FOR THE SNIPER TO KNOW WHAT TO DO WITH THE ACCOUNT.
# for example, for a prename snipe with manual bearer, do BEARERHERE:bearer:prename and DO NOT REPLACE "bearer" or "prename" with anything.

# MOJANG ACCOUNT

EMAIL:PASS
# OR (if you have security questions)
EMAIL:PASS:ANSWER:ANSWER:ANSWER

# Microsoft account
# THIS WILL PROMPT YOU TO LOGIN WHEN AUTHENTICATING, FOLLOW THE INSTRUCTIONS DISPLAYED WHEN YOU ARE PROMPTED.
# there is no password input. it will always be shown as 'oauth2-external' when printing a password.
IDENTIFYING_TEXT:ms
# OR FOR PRENAME
IDENTIFYING_TEXT:prename
# IT IS NOT RECOMMENDED TO USE THE EMAIL:PASS METHOD FOR THIS, SHOWN BELOW.
EMAIL:PASS:ms
# OR FOR PRENAME
EMAIL:PASS:prename

# Manual bearer
BEARERHERE:bearer

# You can add :prename for a prename snipe w/ manual bearer
BEARERHERE:bearer:prename
```
