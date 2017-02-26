# Getting started

Are you interested in the next generation of TsuboneSystem? Why not try on your
computer?

## Building requirements
All of the following softwares:

* [GNU Bash](https://www.gnu.org/software/bash/)
* [GNU Make](https://www.gnu.org/software/make/)
* [Go 1.7](https://golang.org/)
* [OpenJDK 8](http://openjdk.java.net/)
* [MariaDB 10.1](https://mariadb.org/)
* [Node.js](https://nodejs.org/)
* [npm](https://npmjs.com/)
* [procps (ps command)](https://gitlab.com/procps-ng/procps)

Platforms other than Arch Linux are not tested. Fixes for other Linux
distributions, Microsoft-supported Windows and macOS are highly appreciated.

## Client requirements
Either of the following browsers:

* [Google Chrome 53 or any later version](https://www.google.com/chrome/browser/)
* [Microsoft Edge](https://www.microsoft.com/ja-jp/windows/microsoft-edge)
* [Microsoft Internet Explorer 11](https://support.microsoft.com/ja-jp/products/internet-explorer)
* [Mozilla Firefox 45 or any later version](https://www.mozilla.org/en-US/firefox/)

If you use older Internet Explorer, throw it away and go Firefox.

## Testing
1\. Set an appropriate `GOPATH`.

```
export GOPATH=/somewhere/nice
```

2\. Download TsuboneSystem3 with `go get`.

```
go get -d github.com/kagucho/tsubonesystem3/...
cd $GOPATH/src/github.com/kagucho/tsubonesystem3
```

3\. Create a configuration file as `configuration/configuration.go`.

Configure for `tsubonesystem3` command. See `configuration/example` for an
example and detailed explanations for the configuration.

4\. Create a new database.

```
$ mysql
> CREATE DATABASE tsubonesystem;
> exit
$
```

5\. Deploy the testing tables.

```
$ mysql < test.sql
```

5\. Prepare.

```
$ make prepare
```

6\. Make!

```
make
```

make provides other helpful targets; see `Makefile`.

### Testing Database
The testing database is named `test.sql` and it includes records useful for
testing.

The passwords in `members` table is encrypted with `DBPasswordKey` in
`configuration/example/configuration.go`. The raw passwords are `NthPassword`,
where `N` is the member number. For example, the username and the password for
member 1 is `1stDisplayID` and `1stPassword`.

# License
This software is licensed under AGPL-3.0. See `COPYING.TXT`.

# Standards
This application implicitly conforms to the following standards.

* [Accessible Rich Internet Applications (WAI-ARIA) 1.0](https://www.w3.org/TR/2014/REC-wai-aria-20140320/)
* [CSS Snapshot 2017](https://www.w3.org/TR/css-2017/)
* [DOM Standard](https://dom.spec.whatwg.org/)
* [ECMAScript® 2016 Language Specification](http://www.ecma-international.org/ecma-262/7.0/index.html)
* [Encoding Standard](https://encoding.spec.whatwg.org/)
* [The Go Programming Language Specification](https://golang.org/ref/spec)
* [HTML Standard](https://html.spec.whatwg.org/)
* [URL Standard](https://url.spec.whatwg.org/)
* [XMLHttpRequest Standard](https://xhr.spec.whatwg.org/)
