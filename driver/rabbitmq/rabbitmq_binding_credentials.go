package rabbitmq

type RabbitmqBindingCredentials struct {
	Hostname     string `json:"hostname"`
	Host         string `json:"host"`
	VHost        string `json:"vhost"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Uri          string `json:"uri"`
	DashboardUrl string `json:"dashboard_url"`
	Name         string `json:"name"`
	User         string `json:"user"`
	Pass         string `json:"pass"`
}
