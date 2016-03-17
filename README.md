# Falcon App/Scope Launcher

An alternative Ubuntu Touch app & scope launcher

## What's with the name?

The Falcon is a family of rockets developed by SpaceX for carrying payloads into
space. For more information check out this [Wikipedia article](https://en.wikipedia.org/wiki/Falcon_1).

## Build

* Install NPM
    * `sudo apt-get install npm`
* Install click chroot dependencies
    * `sudo click chroot -a armhf -f ubuntu-sdk-15.04 maint apt-get install golang-go golang-go-linux-arm libglib2.0-dev:armhf crossbuild-essential-armhf nodejs npm`
    * `sudo click chroot -a armhf -f ubuntu-sdk-15.04 maint npm install -g gulp`
    * `sudo click chroot -a armhf -f ubuntu-sdk-15.04 maint ln -s `which nodejs` /usr/bin/node`
* Install build dependencies
    * `npm install -g gulp` (May need sudo)
    * `npm install`
* Test locally
    * `gulp run`
* Build click package
    * `gulp build-click`

## Resources

- [Docs for go-unityscopes](https://godoc.org/launchpad.net/go-unityscopes/v2)
- [Go Scope Tutorial](https://developer.ubuntu.com/en/scopes/tutorials/developing-scopes-go/)

## Logo

The Falcon logo is a modified [rocket icon from Game-icons.net](http://game-icons.net/lorc/originals/rocket.html)

## License

Copyright (C) 2016 [Brian Douglass](http://bhdouglass.com/)

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License version 3, as published
by the Free Software Foundation.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranties of MERCHANTABILITY, SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program.  If not, see <http://www.gnu.org/licenses/>.
