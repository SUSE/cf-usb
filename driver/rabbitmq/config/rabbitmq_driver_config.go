package config

type RabbitmqDriverConfig struct {
	DockerEndpoint string `json:"docker_endpoint"`
	DockerImage    string `json:"docker_image"`
	ImageVersion   string `json:"docker_image_version"`
}
