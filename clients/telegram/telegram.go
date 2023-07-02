package telegram

import "net/http"

type Client struct {
	host     string // хость api
	basePath string // префикс с которого начинаются все запросы
	client   http.Client
}

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
func (c *Client) Updates() {

}

// отправка сообщений пользователя
func (c *Client) Send() {

}
