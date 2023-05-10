package vk_client

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type Client struct {
	host     string
	basePath string
	token    string
	client   http.Client
}

const (
	sendMessageMethod = "messages.send"
	version           = "5.131"
	getLongPollServer = "groups.getLongPollServer"
)

func New(host string, token string) *Client {
	return &Client{
		host:   host,
		token:  token,
		client: http.Client{},
	}
}

func (c *Client) SendMessage(userId int, message string, keyboard string) error {
	fmt.Printf("Sending message")

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10000)

	q := url.Values{}
	q.Add("v", version)
	q.Add("user_id", strconv.Itoa(userId))
	q.Add("message", message)
	q.Add("random_id", strconv.Itoa(n))
	q.Add("keyboard", keyboard)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) LongPollServer() ([]Update, error) {
	q := url.Values{}
	q.Add("v", version)
	q.Add("group_id", strconv.Itoa(220396016))

	data, err := c.doRequest(getLongPollServer, q)
	if err != nil {
		return nil, err
	}

	var server LPServerResponse
	err = json.Unmarshal(data, &server)
	if err != nil {
		return nil, err
	}

	var updates []Update

	updates = c.Updates(server.Response.Server, server.Response.Key, server.Response.Ts)

	return updates, err
}

func (c *Client) Updates(server, key, ts string) []Update {
	u := fmt.Sprintf("%s?act=a_check&key=%s&ts=%s&wait=25", server, key, ts)
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Printf("Cant form request:%s", err.Error())
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.client.Do(req)
	if err != nil {
		fmt.Printf("Cant send request:%s", err.Error())
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Cant read responce:%s", err.Error())
	}

	var lPR LongPoolResponse

	if err = json.Unmarshal(body, &lPR); err != nil {
		fmt.Println(err)
	}

	ts = lPR.Ts

	return lPR.Updates
}

func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join("method", method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
