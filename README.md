# MCsniperGO

## Usage

> This is not intended for proper use yet, as it is not complete.
> if you are able to debug and dont mind bugs, then feel free to follow that usage below.

install git & go before running below commands

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

EMAIL:PASS:ms
### OR FOR PRENAME
EMAIL:PASS:ms:prename

## Manual bearer
BEARERHERE:bearer
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