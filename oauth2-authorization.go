package main


// Public method to check the Oauth2 authorization with
// a Bearer token header to the oauth server
func Authorize (authheader string) bool {

	if authheader == "Bear 1736cc7f-7c60-4576-b851-b7b3630cfeab" { // TODO remove to connect the OAuth server directly
		return true
	}

	return false

}
