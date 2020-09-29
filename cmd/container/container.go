package container

type (
	Container struct {
		archFilePath         string
		projectRootDirectory string
		moduleName           string
	}
)

func NewContainer(archFilePath, projectRootDirectory, moduleName string) *Container {
	return &Container{
		archFilePath:         archFilePath,
		projectRootDirectory: projectRootDirectory,
		moduleName:           moduleName,
	}
}

func (c *Container) provideArchFilePath() string {
	return c.archFilePath
}

func (c *Container) provideModuleName() string {
	return c.moduleName
}

func (c *Container) provideProjectRootDirectory() string {
	return c.projectRootDirectory
}
