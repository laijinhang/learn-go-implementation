package auth

import "net/http"

func AddCredentials(req *http.Request) (added bool) {
	host := req.URL.Hostname()

	netrcOnce.Do(readNetrc)
	for _, l := range netrc {
		if l.machine == host {
			req.SetBasicAuth(l.login, l.password)
			return true
		}
	}
	return false
}
