package pass

type PasswordProvider interface {
	GetCurrentPassword() ([]byte, error)
}

type PasswordGenerator interface {
	GeneratePassword(int) ([]string, error)
}
