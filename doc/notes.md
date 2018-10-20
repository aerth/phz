
todo:

  * lang
	(phz) boltdb support, server and cli (key:value get/update)
	(phz) macro/preprocessing support  (htf?)

  * web
	(web) cookies support (maybe: set cookie if any POST)
	(web) form support (template func, bug a)
	(web) set header support (get header works , bug a)
	(web) handle executable files with valid phz shebang (and dont let shebang reach markdown)
	(web) optimize (not yet!)

MVP:

A version 1.0 will have:

	old school php website with /?page=about working
	phpinfo status func {{ phzinfo }}
	forms, headers, upload/download
	cookies, unique request id
	database (for auth, for template data)
	authentication/authorization on a private page
	front-matter support !
