package redisprovisioner

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pivotal-golang/lager"

	dockerclient "github.com/fsouza/go-dockerclient"
)

type RedisServiceProperties struct {
	DockerEndpoint string `json:"docker_endpoint"`
	DockerImage    string `json:"docker_image"`
	ImageVersion   string `json:"docker_image_version"`
}

type RedisProvisioner struct {
	serviceProperties RedisServiceProperties
	client            *dockerclient.Client
	logger            lager.Logger
}

func NewRedisProvisioner(serviceProperties RedisServiceProperties, logger lager.Logger) RedisProvisionerInterface {
	return &RedisProvisioner{
		serviceProperties: serviceProperties,
		logger:            logger,
	}
}

func (provisioner *RedisProvisioner) Init() error {
	var err error

	provisioner.client, err = provisioner.getClient()

	if err != nil {
		return err
	}

	return nil
}

func (provisioner *RedisProvisioner) CreateContainer(containerName string) error {
	err := provisioner.pullImage(provisioner.serviceProperties.DockerImage, provisioner.serviceProperties.ImageVersion)
	if err != nil {
		return err
	}

	pass := generatePassword(12)

	hostConfig := dockerclient.HostConfig{PublishAllPorts: true}
	createOpts := dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image: provisioner.serviceProperties.DockerImage + ":" + provisioner.serviceProperties.ImageVersion,
			Cmd:   []string{"redis-server", fmt.Sprintf("--requirepass %s", pass)},
		},
		HostConfig: &hostConfig,
		Name:       containerName,
	}

	container, err := provisioner.client.CreateContainer(createOpts)
	if err != nil {
		return err
	}

	provisioner.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *RedisProvisioner) DeleteContainer(containerName string) error {

	containerID, err := provisioner.getContainerId(containerName)
	if err != nil {
		return err
	}

	err = provisioner.client.StopContainer(containerID, 5)
	if err != nil {
		return err
	}

	return provisioner.client.RemoveContainer(dockerclient.RemoveContainerOptions{
		ID:    containerID,
		Force: true,
	})
}

func (provisioner *RedisProvisioner) getClient() (*dockerclient.Client, error) {
	client, err := dockerclient.NewClient(provisioner.serviceProperties.DockerEndpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (provisioner *RedisProvisioner) pullImage(imageName, version string) error {
	var buf bytes.Buffer
	pullOpts := dockerclient.PullImageOptions{
		Repository:   imageName,
		Tag:          version,
		OutputStream: &buf,
	}

	err := provisioner.client.PullImage(pullOpts, dockerclient.AuthConfiguration{})
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *RedisProvisioner) findImage(imageName string) (*dockerclient.Image, error) {
	image, err := provisioner.client.InspectImage(imageName)
	if err != nil {
		return nil, fmt.Errorf("Could not find base image %s: %s", imageName, err.Error())
	}

	return image, nil
}

func (provisioner *RedisProvisioner) getContainerId(containerName string) (string, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func (provisioner *RedisProvisioner) getContainer(containerName string) (dockerclient.APIContainers, error) {
	opts := dockerclient.ListContainersOptions{
		All: true,
	}
	containers, err := provisioner.client.ListContainers(opts)
	if err != nil {
		return dockerclient.APIContainers{}, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.TrimPrefix(n, "/") == containerName {
				return c, nil
			}
		}
	}

	return dockerclient.APIContainers{}, fmt.Errorf("Container %s not found", containerName)
}

func (provisioner *RedisProvisioner) GetCredentials(containerName string) (map[string]string, error) {

	m := make(map[string]string)

	container, err := provisioner.getContainer(containerName)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	re := regexp.MustCompile(`'--requirepass\s(\S+)'`)
	submatch := re.FindStringSubmatch(container.Command)
	if submatch == nil {
		return nil, fmt.Errorf("Could not get password")
	}

	m["password"] = submatch[1]
	m["port"] = strconv.FormatInt(container.Ports[0].PublicPort, 10)

	return m, nil
}

func (provisioner *RedisProvisioner) ContainerExists(containerName string) (bool, error) {
	_, err := provisioner.getContainer(containerName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (provisioner *RedisProvisioner) PingServer() error {
	_, err := provisioner.client.Info()
	return err
}

func generatePassword(size int) string {
	dictionary := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	var bytes = make([]byte, size)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}
