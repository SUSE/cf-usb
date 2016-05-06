package csm

import (
	"net/url"
	"strings"
	"testing"

	"github.com/pivotal-golang/lager/lagertest"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("csm-client-test")

var csmEndpoint = "http://192.168.77.77:8081"
var authToken = "csm-auth-token"

func getCSMClient() (CSMInterface, error) {

	csmUrl, err := url.Parse(csmEndpoint)
	if err != nil {
		return nil, err
	}

	client := NewCSMClient(logger, *csmUrl, authToken)
	return client, nil
}

func TestCSMClient(t *testing.T) {
	assert := assert.New(t)
	workspaceID := uuid.NewV4().String()
	//TODO: Remove this line with CSM accepts guids as IDs...
	workspaceID = strings.Replace(workspaceID, "-", "", -1)

	connectionID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = client.CreateWorkspace(workspaceID)
	assert.Nil(err)

	exists, err := client.WorkspaceExists(workspaceID)
	assert.Nil(err)
	assert.True(exists)

	credentials, err := client.CreateConnection(workspaceID, connectionID)
	assert.Nil(err)
	assert.NotNil(credentials)

	credExist, err := client.ConnectionExists(workspaceID, connectionID)
	assert.Nil(err)
	assert.True(credExist)

	client.DeleteConnection(workspaceID, connectionID)
	assert.Nil(err)

	err = client.DeleteWorkspace(workspaceID)
	assert.Nil(err)

}

func TestGetConnectionDoesNotExist(t *testing.T) {
	assert := assert.New(t)
	workspaceID := uuid.NewV4().String()
	//TODO: Remove this line with CSM accepts guids as IDs...
	workspaceID = strings.Replace(workspaceID, "-", "", -1)

	connectionID := uuid.NewV4().String()

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	credExist, err := client.ConnectionExists(workspaceID, connectionID)
	assert.Nil(err)
	assert.False(credExist)
}

func TestGetWorkspaceDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	workspaceID := uuid.NewV4().String()
	//TODO: Remove this line with CSM accepts guids as IDs...
	workspaceID = strings.Replace(workspaceID, "-", "", -1)

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	exists, err := client.WorkspaceExists(workspaceID)
	assert.Nil(err)
	assert.False(exists)
}

func TestDeleteWorkspaceNotExist(t *testing.T) {
	assert := assert.New(t)

	workspaceID := uuid.NewV4().String()
	//TODO: Remove this line with CSM accepts guids as IDs...
	workspaceID = strings.Replace(workspaceID, "-", "", -1)

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = client.DeleteWorkspace(workspaceID)
	//TODO: when is fixed in CSM uncomment the following line
	//assert.NotNil(err)
}

func TestCreateWorkspaceThatExists(t *testing.T) {
	assert := assert.New(t)

	workspaceID := uuid.NewV4().String()
	//TODO: Remove this line with CSM accepts guids as IDs...
	workspaceID = strings.Replace(workspaceID, "-", "", -1)

	client, err := getCSMClient()
	if err != nil {
		assert.Fail(err.Error())
	}

	err = client.CreateWorkspace(workspaceID)
	assert.Nil(err)
	err = client.CreateWorkspace(workspaceID)
	//TODO: uncomment the following line when this is fixed in CSM
	//assert.NotNil(err)
	err = client.DeleteWorkspace(workspaceID)
	assert.Nil(err)
}
