package pass

import "github.com/sethvargo/go-password/password"

type simpleProvider struct {
	passwordLength int
}

func (sp *simpleProvider) GeneratePassword(passwords int) ([]string, error) {
	pswds := make([]string, 0, passwords)
	for i := 0; i < passwords; i++ {
		if pwd, err := password.Generate(sp.passwordLength, 2, 2, false, true); err != nil {
			return nil, err
		} else {
			pswds = append(pswds, pwd)
		}
	}
	return pswds, nil
}

func NewPasswordGenerator(pwdLen int) *simpleProvider {
	return &simpleProvider{
		passwordLength: pwdLen,
	}
}

var _ PasswordGenerator = &simpleProvider{}
