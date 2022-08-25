package decorator

import "github.com/edgarSucre/chat/internal/adapter"

type AdminUseCaseDerator func(adapter.AdminUseCase) adapter.AdminUseCase

func AdminUseCaseWith(uc adapter.AdminUseCase, decs ...AdminUseCaseDerator) adapter.AdminUseCase {
	for _, dec := range decs {
		uc = dec(uc)
	}

	return uc
}
