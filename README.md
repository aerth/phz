# phz

## phz-cli

  * run a single or multiple phz files, print to stdout, debug to stderr

## phz server, phzd

  * serves a directory, on a port

  * any phz files will be parsed, and executed

  * in a phz file you can use markdown and a simple template language (use `<!DOCTYPE html>` to skip markdown processing)

  * there are a number of built-in functions

  * HTML comments and `{{/* multiline comments */}}` are not printed

  * heavily customizable TOML config, with custom error pages


Example .phz syntax:

```
The time is: {{ .Now }}

Working directory: {{ .Env.PWD }}
```

Execute OS commands:

`{{ exec "uname -a" }}` 

Define templates:

`{{define "mytemplate"}}<div>this works</div>{{end}}`

Include templates:

`{{template "mytemplate"}}`


For more see TEMPLATES.md
