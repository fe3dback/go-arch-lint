package operations

import "context"

type (
	repositoryPrivate interface {
		Find1(ctx context.Context)
	}

	RepositoryPublic interface {
		Find2(ctx context.Context)
	}
)

func PublicFindInPrivate2(x repositoryPrivate) {
	_ = x
}

func PublicFindInPublic2(x RepositoryPublic) {
	_ = x
}

func privateFindInPrivate2(x repositoryPrivate) {
	_ = x
}

func privateFindInPublic2(x RepositoryPublic) {
	_ = x
}
