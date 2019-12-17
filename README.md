# zauth

## Build and Run

The build process requires:

- Git
- Go 1.13+
- MySQL 5.7+
- Task
- Minify
- Packr

We will leave installing Go and MySQL to you. Installing Task, Minify, and Packr
can be done like this:

```sh
$ go get -u github.com/go-task/task/cmd/task
$ go get -u github.com/tdewolff/minify/cmd/minify
$ go get -u github.com/gobuffalo/packr/packr
```

After that you can now build and run the application by using Task from the
command line

```sh
$ task build
$ ./zauth
```

Or you can have Task build and run the app for you. This is nice, because it
will recompile and run the new application if you change any source files:

```sh
$ task run
```


## FAQ

### I get "Forbidden - CSRF token invalid" when logging in!

Check your `config.yml` for `production`. If it is set to `true` and you're not
running via SSL/TLS, the CSRF protection will not work. Change to `false` if
you're simply developing locally, or make sure it's served behind SSL.
