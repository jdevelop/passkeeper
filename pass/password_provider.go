package pass

type PasswordProvider interface {
	GetCurrentPassword() ([]byte, error)
}
