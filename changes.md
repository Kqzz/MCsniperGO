I attempted to add a turbo command and also clean up some code / unnessasary things..

some things i suggest implementing or changing.

- your sniper function makes it hard to implement turbo and other commands i wanted to use. the setup itself is okay, but acc types and other info can easily be gained with a helper function that reqs the mojang api.
- because you dont cache bearers and reauth the accs every request, it makes it hard for improving the "system" for the way it was before it was doable, but if you want to add more commands your community wants it will get in the way some.

# Acc type function

```go
func isGC(bearer string) string {
	var accountT string
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)

	fmt.Fprintln(conn, "GET /minecraft/profile/namechange HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Dismal/1.0\r\nAuthorization: Bearer "+bearer+"\r\n\r\n")

	e := make([]byte, 12)
	conn.Read(e)

	switch string(e[9:12]) {
	case `404`:
		accountT = "Giftcard"
	default:
		accountT = "Microsoft"
	}

	return accountT
}
```

# Bearer cache system

```go
func authAccs() {
	var AccountsVer []string
	file, _ := os.Open("accounts.txt")

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		AccountsVer = append(AccountsVer, scanner.Text())
	}

	if len(AccountsVer) == 0 {
		sendE("Unable to continue, you have no accounts added.\n")
		os.Exit(0)
	}

	grabDetails(AccountsVer)

	if !acc.ManualBearer {
		if acc.Bearers == nil {
			sendE("No bearers have been found, please check your details.")
			os.Exit(0)
		} else {
			checkifValid()

			for _, acc := range acc.Bearers {
				if acc.NameChange {
					bearers.Details = append(bearers.Details, apiGO.Info{
						Bearer:      acc.Bearer,
						AccountType: acc.Type,
						Email:       acc.Email,
					})
				}
			}

			if bearers.Details == nil {
				sendE("Failed to authorize your bearers, please rerun the sniper.")
				os.Exit(0)
			}
		}
	}
}

func grabDetails(AccountsVer []string) {
	if acc.ManualBearer {
		for _, bearer := range AccountsVer {
			if apiGO.CheckChange(bearer) {
				bearers.Details = append(bearers.Details, apiGO.Info{
					Bearer:      bearer,
					AccountType: isGC(bearer),
				})
			}

			time.Sleep(time.Second)
		}
	} else {
		if acc.Bearers == nil {
			bearerz := apiGO.Auth(AccountsVer)
			if len(bearerz.Details) == 0 {
				sendE("Unable to authenticate your account(s), please Reverify your login details.\n")
				return
			} else {
				for _, accs := range bearerz.Details {
					acc.Bearers = append(acc.Bearers, apiGO.Bearers{
						Bearer:       accs.Bearer,
						AuthInterval: 86400,
						AuthedAt:     time.Now().Unix(),
						Type:         accs.AccountType,
						Email:        accs.Email,
						Password:     accs.Password,
						NameChange:   apiGO.CheckChange(accs.Bearer),
					})
				}
				acc.SaveConfig()
				acc.LoadState()
			}
		} else {
			if len(acc.Bearers) < len(AccountsVer) {
				var auth []string
				check := make(map[string]bool)

				for _, acc := range acc.Bearers {
					check[acc.Email+":"+acc.Password] = true
				}

				for _, accs := range AccountsVer {
					if !check[accs] {
						auth = append(auth, accs)
					}
				}

				bearerz := apiGO.Auth(auth)

				if len(bearerz.Details) != 0 {
					for _, accs := range bearerz.Details {
						acc.Bearers = append(acc.Bearers, apiGO.Bearers{
							Bearer:       accs.Bearer,
							AuthInterval: 86400,
							AuthedAt:     time.Now().Unix(),
							Type:         accs.AccountType,
							Email:        accs.Email,
							Password:     accs.Password,
							NameChange:   apiGO.CheckChange(accs.Bearer),
						})
					}

					acc.SaveConfig()
					acc.LoadState()
				}
			} else if len(AccountsVer) < len(acc.Bearers) {
				for _, accs := range AccountsVer {
					for _, num := range acc.Bearers {
						if accs == num.Email+":"+num.Password {
							acc.Bearers = append(acc.Bearers, num)
						}
					}
				}
				acc.SaveConfig()
				acc.LoadState()
			}
		}
	}
}

func checkifValid() {
	var reAuth []string
	for _, accs := range acc.Bearers {
		f, _ := http.NewRequest("GET", "https://api.minecraftservices.com/minecraft/profile/name/boom/available", nil)
		f.Header.Set("Authorization", "Bearer "+accs.Bearer)
		j, _ := http.DefaultClient.Do(f)

		if j.StatusCode == 401 {
			sendI(fmt.Sprintf("Account %v turned up invalid. Attempting to Reauth", accs.Email))
			reAuth = append(reAuth, accs.Email+":"+accs.Password)
		}
	}

	if len(reAuth) != 0 {
		sendI(fmt.Sprintf("Reauthing %v accounts..", len(reAuth)))
		bearerz := apiGO.Auth(reAuth)

		if len(bearerz.Details) != 0 {
			for point, data := range acc.Bearers {
				for _, accs := range bearerz.Details {
					if data.Email == accs.Email {
						data.Bearer = accs.Bearer
						data.NameChange = apiGO.CheckChange(accs.Bearer)
						data.Type = accs.AccountType
						data.Password = accs.Password
						data.Email = accs.Email
						data.AuthedAt = time.Now().Unix()
						acc.Bearers[point] = data
						acc.SaveConfig()
					}
				}
			}
		}
	}

	acc.LoadState()
}
```

With this you can replace my api with yours. aswell change the structs to your acc structs.

# Want a reauth system?..

```go
func checkAccs() {
	for {
		time.Sleep(time.Second * 10)

		// check if the last auth was more than a minute ago
		for _, accs := range acc.Bearers {
			if time.Now().Unix() > accs.AuthedAt+accs.AuthInterval {
				sendI(accs.Email + " is due for reauth")

				// authenticating account
				bearers := apiGO.Auth([]string{accs.Email + ":" + accs.Password})

				if bearers.Details != nil {
					for point, data := range acc.Bearers {
						for _, accs := range bearers.Details {
							if data.Email == accs.Email {
								data.Bearer = accs.Bearer
								data.NameChange = apiGO.CheckChange(accs.Bearer)
								data.Type = accs.AccountType
								data.Password = accs.Password
								data.Email = accs.Email
								data.AuthedAt = time.Now().Unix()
								acc.Bearers[point] = data
							}
						}
					}

					acc.SaveConfig()
					acc.LoadState()
					break // break the loop to update the info.
				}

				// if the account isnt usable, remove it from the list
				var ts apiGO.Config
				for _, i := range acc.Bearers {
					if i.Email != accs.Email {
						ts.Bearers = append(ts.Bearers, i)
					}
				}

				acc.Bearers = ts.Bearers

				acc.SaveConfig()
				acc.LoadState()
				break // break the loop to update the info.
			}
		}
	}
}
```

this is code taken and modified from me and peets sniper alien..

now you DONT need to make a system like this but its pretty easy to implement.. and will help you ALOT.

