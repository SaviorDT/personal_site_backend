package apipaths

const (
	AuthGroup = "/auth"

	GitHubLoginRel    = "/login-github"
	GitHubCallbackRel = "/login-github-callback"
	GoogleLoginRel    = "/login-google"
	GoogleCallbackRel = "/login-google-callback"

	GitHubLoginPath    = AuthGroup + GitHubLoginRel
	GitHubCallbackPath = AuthGroup + GitHubCallbackRel
	GoogleLoginPath    = AuthGroup + GoogleLoginRel
	GoogleCallbackPath = AuthGroup + GoogleCallbackRel
)
