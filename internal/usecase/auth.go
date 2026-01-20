package usecase

type CheckAuth struct {
	checker AuthChecker
}

type AuthChecker interface {
	IsUserAllowed(userID int64) bool
}

func NewCheckAuth(checker AuthChecker) *CheckAuth {
	return &CheckAuth{
		checker: checker,
	}
}

func (a *CheckAuth) Execute(userID int64) bool {
	return a.checker.IsUserAllowed(userID)
}
