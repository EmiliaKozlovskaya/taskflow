package web_fs_repository

// здесь Repository ни от чего зависеть не будет, потому что чтобы с файловой системы прочитать файл никакого подключения или пула не нужно
type WebRepository struct {
}

func NewWebRepository() *WebRepository {
	return &WebRepository{}
}
