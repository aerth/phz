package phz

func handleaqua(s string) string {
	return execslice([]string{"aquachain", "attach", "--exec", s})
}
