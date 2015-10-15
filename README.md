# Trending Machine

Explore Github's trending page in the past. You can see it running here: [http://trending-machine.salimalami.com](http://trending-machine.salimalami.com).

This project is intended to be an example of the [spider package](https://github.com/celrenheit/spider).

Is has a list of pre-defined languages defined in the [settings.json](https://github.com/celrenheit/trending-machine/blob/master/settings.json).
It could be done for all languages but it will put too much load on the Github's servers.

# Installation

You need [Go](https://golang.org/), [Node.js](https://nodejs.org/), [NPM](https://www.npmjs.com/), [Gulp](http://gulpjs.com/) and [Bower](http://bower.io/) installed on your computer.
For the database you'll need [MongoDB](http://www.mongodb.com/).

```shell
$ go get -u github.com/celrenheit/trending-machine
```

To install the dependencies of the web app you need to run:
```shell
$ cd $GOPATH/src/github.com/celrenheit/trending-machine/web
$ npm i && bower i && gulp build --prod
```

# Usage

```shell
$ cd $GOPATH/src/github.com/celrenheit/trending-machine
$ go run *.go
```

This will launch a web server. Open a new tab/window in your browser pointing to [http://localhost:3000](http://localhost:3000).

# Developpement

Open two terminal windows.
In the first one run the following:

```shell
$ cd $GOPATH/src/github.com/celrenheit/trending-machine
$ go run *.go
```

And in the second run:

```shell
$ cd $GOPATH/src/github.com/celrenheit/trending-machine/web
$ gulp
```


# Inspiration

This project is inspired by [https://github.com/josephyzhou/github-trending](https://github.com/josephyzhou/github-trending).
