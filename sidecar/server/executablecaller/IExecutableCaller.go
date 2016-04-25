package executablecaller

type IExecutableCaller interface {
	GetCreateWorkspaceExecutable() string

	CreateWorkspaceCaller(arg1 string) (string, int)
	GetCreateConnectionExecutable() string
	CreateConnectionCaller(arg1 string, arg2 string) (string, int)
	GetDeleteWorkspaceExecutable() string
	DeleteWorkspaceCaller(arg1 string) (string, int)
	GetDeleteConnectionExecutable() string

	DeleteConnectionCaller(arg1 string, arg2 string) (string, int)

	GetGetConnectionExecutable() string
	GetConnectionCaller(arg1 string, arg2 string) (string, int)

	GetGetWorkspaceExecutable() string

	GetWorkspaceCaller(arg1 string) (string, int)
}
