# zauth

Lightweight authentication server, with basic LDAP support, and a web UI.

## Compile

The build process requires:

- Git
- Go 1.13+
- MySQL 5.7+
- Task
- Goimports
- Minify
- Packr

We will leave installing Go and MySQL to you. Installing Task, Minify, and Packr
can be done like this:

```sh
go get -u github.com/go-task/task/v2/cmd/task
go get -u golang.org/x/tools/cmd/goimports
go get -u github.com/tdewolff/minify/cmd/minify
go get -u github.com/gobuffalo/packr/packr
```

After that you can now build and run the application by using Task from the
command line

```sh
task build
./zauth
```

Or you can have Task build and run the app for you. This is nice, because it
will recompile and run the new application if you change any source files:

```sh
task run
```

## FAQ

### I get "Forbidden - CSRF token invalid" when logging in!

Check your `config.yml` for `production`. If it is set to `true` and you're not
running via SSL/TLS, the CSRF protection will not work. Change to `false` if
you're simply developing locally, or make sure it's served behind SSL.

### How can I query and test the LDAP server?

One way is to install `ldapsearch` which is standards compliant. Then you can:

```sh
# Get info about the user joshz (anonymously)
$ ldapsearch -x -H ldap://localhost:3389 "uid=joshz"
# Get info about the group admin (anonymously)
$ ldapsearch -x -H ldap://localhost:3389 "cn=admin"
# ldapwhoami currently doesn't work with zauth (BUG), so you can use to login
$ ldapsearch -x -H ldap://localhost:3389 -W -D 'uid=joshz'
#
# Login as joshz (and return all of your info) (currently broken)
$ ldapwhoami -x -H ldap://localhost:3389 -W -D 'uid=joshz'

```
