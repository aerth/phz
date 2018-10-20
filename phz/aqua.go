package phz

func init() {
	// load plugin
	DefaultFuncMap["aqua"] = handleaqua
}

// aquachain js exec wrapper (connects via ipc socket)
func handleaqua(s string) string {
	return execslice([]string{"aquachain", "attach", "--exec", s})
}
