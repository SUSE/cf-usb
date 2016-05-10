package csm

type CSMInterface interface {
	Login(string, string) error
	CreateWorkspace(string) error
	WorkspaceExists(string) (bool, error)
	DeleteWorkspace(string) error
	CreateConnection(string, string) (interface{}, error)
	ConnectionExists(string, string) (bool, error)
	DeleteConnection(string, string) error
}
