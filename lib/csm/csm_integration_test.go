package csm

import (
	"os"
	"testing"
	"time"

	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var logger = lagertest.NewTestLogger("csm-client-test")

var csmEndpoint string
var authToken string

func getCSMClient() (CSM, error) {
	client := NewCSMClient(logger)
	//skipping SSL validation
	err := client.Login(csmEndpoint, authToken, "", true)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestCSMClient(t *testing.T) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" || authToken == "" {
		t.Skipf("Skipping test TestCSMClient - missing CSM_ENDPOINT and CSM_API_KEY")
	}
	assert := assert.New(t)
	workspaceID := uuid.NewV4().String()

	connectionID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = client.CreateWorkspace(workspaceID)
	time.Sleep(120 * time.Second)
	assert.Nil(err)

	exists, isNoop, err := client.WorkspaceExists(workspaceID)
	assert.Nil(err)
	if isNoop == false {
		assert.True(exists)
	}

	credentials, err := client.CreateConnection(workspaceID, connectionID)
	time.Sleep(120 * time.Second)
	assert.Nil(err)
	assert.NotNil(credentials)

	credExist, isNoop, err := client.ConnectionExists(workspaceID, connectionID)
	assert.Nil(err)
	if isNoop == false {
		assert.True(credExist)
	}

	client.DeleteConnection(workspaceID, connectionID)
	time.Sleep(120 * time.Second)
	assert.Nil(err)

	err = client.DeleteWorkspace(workspaceID)
	time.Sleep(120 * time.Second)
	assert.Nil(err)

}

func TestGetConnectionDoesNotExist(t *testing.T) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" || authToken == "" {
		t.Skipf("Skipping test TestCSMClient - missing CSM_ENDPOINT and CSM_API_KEY")
	}
	assert := assert.New(t)
	workspaceID := uuid.NewV4().String()

	connectionID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	credExist, isNoop, err := client.ConnectionExists(workspaceID, connectionID)
	assert.Nil(err)
	if isNoop == false {
		assert.False(credExist)
	}
}

func TestGetWorkspaceDoesNotExist(t *testing.T) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" || authToken == "" {
		t.Skipf("Skipping test TestCSMClient - missing CSM_ENDPOINT and CSM_API_KEY")
	}
	assert := assert.New(t)

	workspaceID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	exists, isNoop, err := client.WorkspaceExists(workspaceID)
	assert.Nil(err)
	if isNoop == false {
		assert.False(exists)
	}
}

func TestDeleteWorkspaceNotExist(t *testing.T) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" || authToken == "" {
		t.Skipf("Skipping test TestCSMClient - missing CSM_ENDPOINT and CSM_API_KEY")
	}
	assert := assert.New(t)

	workspaceID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = client.DeleteWorkspace(workspaceID)
	//TODO: when is fixed in CSM uncomment the following line
	//assert.NotNil(err)
}

func TestCreateWorkspaceThatExists(t *testing.T) {
	csmEndpoint = os.Getenv("CSM_ENDPOINT")
	authToken = os.Getenv("CSM_API_KEY")
	if csmEndpoint == "" || authToken == "" {
		t.Skipf("Skipping test TestCSMClient - missing CSM_ENDPOINT and CSM_API_KEY")
	}
	assert := assert.New(t)

	workspaceID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = client.CreateWorkspace(workspaceID)
	assert.Nil(err)
	err = client.CreateWorkspace(workspaceID)
	assert.NotNil(err)
	err = client.DeleteWorkspace(workspaceID)
	assert.Nil(err)
}
