package telegram

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"proj/lib/e"
	"strconv"
)

type Client struct {
	host     string // хост api
	basePath string // префикс с которого начинаются все запросы
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}
func newBasePath(token string) string {
	return "bot" + token // чисто штук
}

// получение новых сообщений
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}
	var res UpdateResponse
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}
	return res.Result, nil
}

// отправка сообщений пользователя
func (c *Client) Send(chatId int, text string) error {

	q := url.Values{}
	q.Add("chatId", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}
	return nil
}
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const errMsg = "can't do request"
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	return body, nil
}
