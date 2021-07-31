package blizz

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/levigross/grequests"
)

type BlizzardToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *client) BlizzAuthRoutine() {
	delay := time.Duration(c.cfg.AuthTimeOut) * time.Hour
	t := time.NewTicker(delay)

	defer t.Stop()

	for {
		select {
		case <-t.C:
			c.log.Info("Making blizzard auth call")

			if err := c.MakeBlizzAuth(); err != nil {
				c.log.Errorf("error making blizzard auth request: %v", err)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *client) MakeBlizzAuth() error {
	body := strings.NewReader("grant_type=client_credentials")
	ro := &grequests.RequestOptions{
		RequestBody: body,
		Auth:        []string{c.cfg.ClientID, c.cfg.ClientSecret},
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	response, err := c.session.Post(c.cfg.AUTHUrl, ro)
	if err != nil {
		return errors.Wrapf(err, "MakeBlizzAuth POST")
	}

	tokenData := new(BlizzardToken)
	if err := response.JSON(tokenData); err != nil {
		return errors.Wrapf(err, "MakeBlizzAuth Parse JSON")
	}

	c.token = tokenData

	return nil
}
