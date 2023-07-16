package di

import (
	"github.com/fe3dback/go-arch-lint/internal/pkg/deepscan/test/project/internal/operations"
	"github.com/fe3dback/go-arch-lint/internal/pkg/deepscan/test/project/internal/repository"
)

func TestCases() {
	t1c1DirectPass()
	t1c2NestedVariablePass()
	t1c3FunctionPass()
	t1c4VariablePass()
	t1c5InterfacePass()
	t1c6DualInterfacePass()
	t1c7DualActualInterfacePass()
	t1c8BasicSpread()
	t1c9Slice()
	t1c10Array()
	t1c11SpreadNames()
	t1c12SpreadAnon()
}

func t1c1DirectPass() {
	operations.NewProcessorBasic1(
		repository.NewMemory(),
	)
}

func t1c2NestedVariablePass() {
	a := repository.NewMemory()
	b := a
	c := b
	d := c

	operation := operations.NewProcessorNames1(d)
	_ = operation
}

func t1c3FunctionPass() {
	operations.NewProcessorAlias1(t4FunctionPassProvideMemoryRepository())
}

func t1c4VariablePass() {
	memRepo := repository.NewMemory()
	anotherMemRepo := repository.NewMemory()
	operation := operations.NewProcessorBasicDual1(memRepo, anotherMemRepo)
	_ = operation
}

func t4FunctionPassProvideMemoryRepository() *repository.Memory {
	return repository.NewMemory()
}

func t1c5InterfacePass() {
	operations.NewProcessorBasic1(t5InterfaceProvide())
}

func t5InterfaceProvide() operations.PublicFetcherForDI {
	return repository.NewMemory()
}

func t1c6DualInterfacePass() {
	operations.NewProcessorBasic1(t6DualInterfaceProvideA())
}

func t6DualInterfaceProvideA() operations.PublicFetcherForDI {
	return t6DualInterfaceProvideB()
}

func t6DualInterfaceProvideB() operations.PublicFetcherForDI {
	return repository.NewMemory()
}

func t1c7DualActualInterfacePass() {
	operations.NewProcessorBasic1(t7DualActualInterfaceProvideA())
}

func t7DualActualInterfaceProvideA() operations.PublicFetcherForDI {
	return t7DualActualInterfaceProvideB()
}

func t7DualActualInterfaceProvideB() *repository.Memory {
	return repository.NewMemory()
}

func t1c8BasicSpread() {
	sameRepo := repository.NewMemory()

	operations.NewProcessorBasicSpreadTypes1(
		repository.NewMemory(),
		sameRepo,
		func() *repository.Memory {
			return repository.NewMemory()
		}(),
	)
}

func t1c9Slice() {
	repo1 := repository.NewMemory()
	repo2 := repository.NewMemory()

	operations.NewProcessorBasicSlice1([]operations.PublicFetcherForDI{
		repo1,
		repo2,
		repo1,
	})
}

func t1c10Array() {
	repo1 := repository.NewMemory()
	repo2 := repository.NewMemory()

	operations.NewProcessorBasicArray1([5]operations.PublicFetcherForDI{
		repo1, repo2, repo1,
	})
}

func t1c11SpreadNames() {
	repo1 := repository.NewMemory()
	operations.NewProcessorBasicSpreadNames1(repository.NewMemory(), repo1)
}

func t1c11SpreadNamesCopy() {
	repo2 := repository.NewMemory()
	operations.NewProcessorBasicSpreadNames1(repo2, repository.NewMemory())
}

func t1c12SpreadAnon() {
	repo1 := repository.NewMemory()
	repo2 := repository.NewMemory()
	operations.NewProcessorBasicSpreadNamesAnonim1(repo1, repo2, repo1, repo2)
}
