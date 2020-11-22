package multimc

import (
	"cornstone/util"
)

type AccountsJson struct {
	Accounts      []Accounts `json:"accounts"`
	ActiveAccount string     `json:"activeAccount"`
	FormatVersion int        `json:"formatVersion"`
}
type Profiles struct {
	ID     string `json:"id"`
	Legacy bool   `json:"legacy"`
	Name   string `json:"name"`
}
type User struct {
	ID string `json:"id"`
}
type Accounts struct {
	AccessToken   string     `json:"accessToken"`
	ActiveProfile string     `json:"activeProfile"`
	ClientToken   string     `json:"clientToken"`
	Profiles      []Profiles `json:"profiles"`
	User          User       `json:"user"`
	Username      string     `json:"username"`
}

func MakeNewAccountsJson(name string) ([]byte, error) {
	hash := "ba1f2511fc30423bdbb183fe33f3dd0f"
	email := "test@account.dev"
	accounts := AccountsJson{
		Accounts: []Accounts{
			{
				AccessToken:   hash,
				ActiveProfile: hash,
				ClientToken:   hash,
				Profiles: []Profiles{
					{
						ID:     hash,
						Legacy: false,
						Name:   name,
					},
				},
				User:     User{ID: hash},
				Username: email},
		},
		ActiveAccount: email,
		FormatVersion: 2,
	}
	result, err := util.JsonMarshalPretty(accounts)
	if err != nil {
		return nil, err
	}
	return result, nil
}
