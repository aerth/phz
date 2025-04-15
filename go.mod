module x/phzd

go 1.24.2

// +heroku install ./cmd/phzd

require (
	github.com/BurntSushi/toml v1.5.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/russross/blackfriday v2.0.0+incompatible
	golang.org/x/crypto v0.37.0
)

require (
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)
