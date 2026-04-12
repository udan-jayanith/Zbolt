package url_utils

import "net/url"

func CleanURL(u *url.URL) {
	u.RawQuery = ""
	u.RawFragment = ""
	u.Fragment = ""
	u.ForceQuery = false
}
