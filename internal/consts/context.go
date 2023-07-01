package consts

type contextIntKey int
type contextStringKey string

const (
	// CtxUserCode const
	CtxUserCode contextIntKey = iota
	// CtxUserEmail const
	CtxUserEmail
	// CtxUserPhone const
	CtxUserPhone
	// CtxIP const
	CtxIP
	// CtxUserAgent const
	CtxUserAgent
	// CtxLang const
	CtxLang
	// CtxUserInfo const
	CtxUserInfo
	// CtxAuthInfo auth context key
	CtxAuthInfo
)

const (
	// CtxAccess const
	CtxAccess contextStringKey = "access"
)
