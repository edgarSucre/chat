package usecase

type AdminUsecaseOption func(*AdminUsecase)

func WithHasher(h Secure) AdminUsecaseOption {
	return func(uc *AdminUsecase) {
		uc.hasher = h
	}
}
