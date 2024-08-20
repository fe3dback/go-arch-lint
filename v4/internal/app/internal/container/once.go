package container

import (
	"fmt"
	"runtime"

	"github.com/fe3dback/go-arch-lint-sdk/pkg/safemap"
)

var instancesCache = safemap.New[string, any]()

// once выполняет вызов функции ровно один раз, учитывая конкретное место вызова once (путь к файлу и номер строки).
// При последующих вызовах в этом же месте возвращает кэшированное значение
func once[T any](factory func() T) T {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("failed to build caller key")
	}

	key := fmt.Sprintf("%s:%v", file, line)

	cachedInstance, alreadyCached := instancesCache.Get(key)

	if alreadyCached {
		if typedInstance, ok := cachedInstance.(T); ok {
			return typedInstance
		}

		panic(fmt.Sprintf(
			"failed to cast cached instance to target type (key: %s, expected %T, actual %T)",
			key,
			*new(T), // костыль, чтобы можно было указать ожидаемый тип
			cachedInstance,
		))
	}

	newInstance := factory()
	instancesCache.Set(key, newInstance)

	return newInstance
}
