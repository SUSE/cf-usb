package rabbitmqprovisioner

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	"github.com/hpcloud/cf-usb/driver/rabbitmq/config"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-golang/lager"

	dockerclient "github.com/fsouza/go-dockerclient"
)

const CONTAINER_START_TIMEOUT int = 30

type RabbitmqProvisioner struct {
	driverConfig config.RabbitmqDriverConfig
	client       *dockerclient.Client
	logger       lager.Logger
}

func NewRabbitmqProvisioner(logger lager.Logger) RabbitmqProvisionerInterface {
	return &RabbitmqProvisioner{logger: logger}
}

func (provisioner *RabbitmqProvisioner) Connect(driverConfig config.RabbitmqDriverConfig) error {
	var err error

	dockerUrl, err := url.Parse(driverConfig.DockerEndpoint)
	if err != nil {
		return err
	}

	if dockerUrl.Scheme == "" {
		return errors.New("Invalid URL format")
	}

	provisioner.driverConfig = driverConfig
	provisioner.client, err = provisioner.getClient()

	if err != nil {
		return err
	}

	return nil
}

func (provisioner *RabbitmqProvisioner) CreateContainer(containerName string) error {
	err := provisioner.pullImage(provisioner.driverConfig.DockerImage, provisioner.driverConfig.ImageVersion)
	if err != nil {
		return err
	}

	admin_user, err := secureRandomString(32)
	if err != nil {
		return err
	}
	admin_pass, err := secureRandomString(32)
	if err != nil {
		return err
	}
	hostConfig := dockerclient.HostConfig{PublishAllPorts: true}
	createOpts := dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image: provisioner.driverConfig.DockerImage + ":" + provisioner.driverConfig.ImageVersion,
			Env: []string{"RABBITMQ_DEFAULT_USER=" + admin_user,
				"RABBITMQ_DEFAULT_PASS=" + admin_pass},
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

	retry := 1
	for retry < CONTAINER_START_TIMEOUT {
		state, err := provisioner.getContainerState(containerName)
		if err != nil {
			return err
		}
		if state.Running {
			break
		}
		retry++
	}

	return nil
}

func (provisioner *RabbitmqProvisioner) DeleteContainer(containerName string) error {

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

func (provisioner *RabbitmqProvisioner) getClient() (*dockerclient.Client, error) {
	client, err := dockerclient.NewClient(provisioner.driverConfig.DockerEndpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (provisioner *RabbitmqProvisioner) pullImage(imageName, version string) error {
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

func (provisioner *RabbitmqProvisioner) findImage(imageName string) (*dockerclient.Image, error) {
	image, err := provisioner.client.InspectImage(imageName)
	if err != nil {
		return nil, fmt.Errorf("Could not find base image %s: %s", imageName, err.Error())
	}

	return image, nil
}

func (provisioner *RabbitmqProvisioner) getContainerId(containerName string) (string, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func (provisioner *RabbitmqProvisioner) getContainer(containerName string) (dockerclient.APIContainers, error) {
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

func (provisioner *RabbitmqProvisioner) inspectContainer(containerId string) (*dockerclient.Container, error) {
	return provisioner.client.InspectContainer(containerId)
}

func (provisioner *RabbitmqProvisioner) getAdminCredentials(containerName string) (map[string]string, error) {

	m := make(map[string]string)
	containerId, err := provisioner.getContainerId(containerName)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	container, err := provisioner.inspectContainer(containerId)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	var env dockerclient.Env
	env = make([]string, len(container.Config.Env)) // container.Config.Env.(dockerclient.Env)  // dockerclient.Env( []string{ container.Config.Env })
	copy(env, container.Config.Env)
	m["user"] = env.Get("RABBITMQ_DEFAULT_USER")
	m["password"] = env.Get("RABBITMQ_DEFAULT_PASS")
	for k, v := range container.NetworkSettings.Ports {
		if k == "15672/tcp" {
			m["mgmt_port"] = v[0].HostPort
		}
		if k == "5672/tcp" {
			m["port"] = v[0].HostPort
		}
	}
	return m, nil
}

func (provisioner *RabbitmqProvisioner) ContainerExists(containerName string) (bool, error) {
	_, err := provisioner.getContainer(containerName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (provisioner *RabbitmqProvisioner) PingServer() error {
	_, err := provisioner.client.Info()
	return err
}

func (provisioner *RabbitmqProvisioner) DeleteUser(containerName, credentialId string) error {
	host, err := provisioner.getHost()
	if err != nil {
		return err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return err
	}
	user, err := getMD5Hash(credentialId)
	if err != nil {
		return err
	}
	_, err = rmqc.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *RabbitmqProvisioner) UserExists(containerName, credentialId string) (bool, error) {
	host, err := provisioner.getHost()
	if err != nil {
		return false, err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return false, err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return false, err
	}
	user, err := getMD5Hash(credentialId)
	if err != nil {
		return false, err
	}
	users, err := rmqc.ListUsers()
	if err != nil {
		return false, err
	}
	if users == nil {
		return false, err
	}
	
	for _, u := range users {
		if u.Name == user{
			return true, nil
		}
	}

	return false, nil
}

func (provisioner *RabbitmqProvisioner) getContainerState(containerName string) (dockerclient.State, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return dockerclient.State{}, nil
	}

	c, err := provisioner.inspectContainer(container.ID)
	if err != nil {
		return dockerclient.State{}, nil
	}
	return c.State, nil
}

func (provisioner *RabbitmqProvisioner) CreateUser(containerName, credentialId string) (map[string]string, error) {
	host, err := provisioner.getHost()
	if err != nil {
		return nil, err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return nil, err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)
	newUser, err := getMD5Hash(credentialId)
	if err != nil {
		return nil, err
	}

	userPass, err := secureRandomString(32)
	if err != nil {
		return nil, err
	}
	_, err = rmqc.PutUser(newUser, rabbithole.UserSettings{Password: userPass, Tags: "management,policymaker"})
	if err != nil {
		return nil, err
	}

	_, err = rmqc.UpdatePermissionsIn("/", newUser, rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"})
	if err != nil {
		return nil, err
	}
	m["host"] = host
	m["user"] = newUser
	m["password"] = userPass
	m["mgmt_port"] = admin["mgmt_port"]
	m["port"] = admin["port"]
	x, err := rmqc.GetVhost("/")
	if err != nil {
		return nil, err
	}
	m["vhost"] = x.Name

	return m, nil
}

func (provisioner *RabbitmqProvisioner) getHost() (string, error) {
	host := ""
	dockerUrl, err := url.Parse(provisioner.driverConfig.DockerEndpoint)
	if err != nil {
		return "", err
	}

	if dockerUrl.Scheme == "unix" {
		host, err = getLocalIP()
		if err != nil {
			return "", err
		}
	} else {
		host = strings.Split(dockerUrl.Host, ":")[0]
	}

	return host, nil
}

func secureRandomString(bytesOfEntpry int) (string, error) {
	rb := make([]byte, bytesOfEntpry)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(rb), nil
}

func getMD5Hash(text string) (string, error) {
	hasher := md5.New()
	hasher.Write([]byte(text))
	generated := hex.EncodeToString(hasher.Sum(nil))

	reg := regexp.MustCompile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(generated, ""), nil
}

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, inface := range interfaces {
		addresses, err := inface.Addrs()
		if err != nil {
			return "", err
		}
		for _, address := range addresses {
			ipnet, ok := address.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipnet.IP.To4()
			if ip == nil || ip.IsLoopback() {
				continue
			}
			return ip.String(), nil
		}
	}
	return "", fmt.Errorf("Could not find IP address")
}
