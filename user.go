package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type userResult struct {
	User User `json:"user"`
}

type usersResult struct {
	Users []User `json:"users"`
}

type User struct {
	Id           int            `json:"id"`
	Login        string         `json:"login"`
	Firstname    string         `json:"firstname"`
	Lastname     string         `json:"lastname"`
	Mail         string         `json:"mail"`
	CreatedOn    string         `json:"created_on"`
	LatLoginOn   string         `json:"last_login_on"`
	Memberships  []Membership   `json:"memberships"`
	CustomFields []*CustomField `json:"custom_fields,omitempty"`
}

func (c *Client) Users() ([]User, error) {
	res, err := c.Get(c.endpoint + "/users.json?key=" + c.apikey + c.getPaginationClause())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r usersResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.Users, nil
}

func (c *Client) LockedUsers() ([]User, error) {
	res, err := c.Get(c.endpoint + "/users.json?key=" + c.apikey + c.getPaginationClause() + "&status=3")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r usersResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return r.Users, nil
}

func (c *Client) User(id int) (*User, error) {
	res, err := c.Get(c.endpoint + "/users/" + strconv.Itoa(id) + ".json?key=" + c.apikey)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var r userResult
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		err = decoder.Decode(&r)
	}
	if err != nil {
		return nil, err
	}
	return &r.User, nil
}

func (c *Client) DeleteUser(id int) error {
	req, err := http.NewRequest("DELETE", c.endpoint+"/users/"+strconv.Itoa(id)+".json?key="+c.apikey, strings.NewReader(""))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return errors.New("Not Found")
	}

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		var er errorsResult
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	}
	return err
}
