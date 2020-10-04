package container

type (
	Container struct {
		archFilePath         string
		projectRootDirectory string
		moduleName           string
		useColors            bool
	}
)

func NewContainer(
	archFilePath string,
	projectRootDirectory string,
	moduleName string,
	useColors bool,
) *Container {
	return &Container{
		archFilePath:         archFilePath,
		projectRootDirectory: projectRootDirectory,
		moduleName:           moduleName,
		useColors:            useColors,
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

func (c *Container) provideFlagUseColors() bool {
	return c.useColors
}
