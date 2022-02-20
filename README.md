# MCsniperGO
<a href="https://mcsniperpy.com/discord"><img src="https://invidget.switchblade.xyz/yp69ZqtxNk"/><a/>

This project was made possible by my [donators](https://kqzz.me/donators)

## Installation

### Windows

- Download the latest EXE from [releases](https://github.com/Kqzz/MCsniperGO/releases/)
- Make a folder for sniping, move the exe into that folder
- Run the exe by double clicking (windows defender may popup, that's normal) 
- Open accounts.txt and **ADD YOUR ACCOUNTS** according to the [FORMAT HERE](https://github.com/Kqzz/MCsniperGO#accounts-file-format)
- Run the sniper by double clicking
- Enter the *dropping* username you want to snipe and your offset
- Wait for the drop! good luck.

### Linux
> *this guide is very brief and assumes knowledge of linux*

- `wget` the binary link (`wget https://github-releases-url-here`) found in [releases](https://github.com/Kqzz/MCsniperGO/releases/)
- Run `chmod +x ./binary_name_here`
- Run `./binary_name_here`
- Open accounts.txt (`nano accounts.txt`) and add your accounts [according the the format below](https://github.com/Kqzz/MCsniperGO#accounts-file-format)
- Run `./binary_name_here`. Note: you need to leave the terminal session open for the sniper to keep running, use `tmux`, `./binary_name_here`, and then `<Ctrl>-<B> <D>` to run in tmux terminal session.
- Wait for the drop! good luck.


### Mac

Download the mac binary, and then follow the Linux settings (ignore the wget step) using `terminal`. Make sure your terminal is in the same folder as the binary you downloaded.

### Compiling from source (not recommended)

Install [git](https://git-scm.com/) and [go](https://golang.org/dl/) 1.17 or later.
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
