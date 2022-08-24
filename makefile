mock-repo:
	mockgen -package mockrepo -destination internal/mock/repo/repo.go github.com/edgarSucre/chat/internal/usecase AdminRepository

mock-hasher:
	mockgen -package mockhash -destination internal/mock/hasher/hasher.go github.com/edgarSucre/chat/internal/usecase Secure