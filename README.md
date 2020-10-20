# notes - a journal server

- ![CI](https://github.com/klingtnet/notes/workflows/CI/badge.svg)
- [Releases][releases]

`notes` is a web service equivalent to a notebook, you can store ideas and look for them with a full text search.  The search makes it even better than a notebook.
Notes are stored in a encrypted SQLite database using SQLCipher.
It supports github flavored markdown and everything is contained in a single binary (execpt the database file).

I am happy to get pull requests and here a couple of improvement ideas:

- [ ] mark tasks as finished by clicking on a checkbox
- [ ] only show open tasks
- [ ] add pagination or date filtering

## Installation

For Linux users with a working Go setup the easiest way is to install the server and its systemd user service using by running `make install`.
You can also [grab a prebuilt binary from the releases page][releases] and copy it into your `$PATH`.
A database file is created on the first start at `$XDG_CONFIG_HOME/notes/notes.db` or wherever you set the argument of `--database-path` to.
It is important to set a proper database password on the first run using `DATABASE_PASSPHRASE` environment variable (you can use the flag as well but [this comes with some security downsides][so-password].)
At the moment there is no option to change the password but you should be able to achieve that in a db shell session (`scripts/dbshell` shows how to start such a session).

[releases]: https://github.com/klingtnet/notes/releases
[so-password]: https://security.stackexchange.com/questions/14000/environment-variable-accessibility-in-linux/14009#14009