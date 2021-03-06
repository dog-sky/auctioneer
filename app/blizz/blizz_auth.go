package blizz

import (
	"fmt"
	"strings"
	"time"

	"github.com/levigross/grequests"
)

func (c *client) BlizzAuthRoutine() {
	delay := 6 * time.Hour
	t := time.NewTicker(delay)
	defer t.Stop()

	for range t.C {
		c.log.Info("Making blizzard auth call")
		if err := c.MakeBlizzAuth(); err != nil {
			c.log.Errorf("error making blizzard auth request: %v", err)
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
		return fmt.Errorf("error making blizzard auth request: %v", err)
	}

	tokenData := new(BlizzardToken)
	if err := response.JSON(tokenData); err != nil {
		return fmt.Errorf("error unmarshaling blizzard auth response: %v", err)
	}

	c.token = tokenData

	return nil
}
