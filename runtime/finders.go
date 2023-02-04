package runtime

func IsOpenRedirectParam(key string) bool {
	var needles []string
	needles = append(needles, "f", "forward", "d", "dest", "destination", "r", "redir", "redirect", "redirector", "u", "uri", "url", "p", "path", "continue", "window", "windows", "to", "out")
	needles = append(needles, "view", "dir", "show", "navigation", "navigate", "navigator", "open", "opener", "file", "val", "validate", "validator", "domain", "call", "back", "callback")
	needles = append(needles, "ret", "return", "returns", "page", "feed", "host", "port", "n", "next", "data", "ref", "reference", "site", "html", "link", "address", "output")

	for _, needle := range needles {
		if needle == key {
			return true
		}
	}

	return false
}
