littlereader
------------

`littlereader` is a minimalistic single-user RSS reader.

It's written in Go, and is fairly simple to get up and running.

1. Clone the repo.
2. Copy your `subscriptions.xml` from your Google Reader takeout.
3. Run `make import`
4. Run `make run`
5. Go to http://localhost:9090/

Features
========

* Support Atom and RSS
* Automatic fetching every 12 hours of new articles
* Mark items as read
* Mark all as read
* Add feeds

TODO
====

* Remove feeds easily
* Any theming

License
-------

All code in this repository, unless otherwise specified, is hereby
licensed under the MIT Public License:

Copyright (c) 2013 Christopher Swenson.

	Permission is hereby granted, free of charge, to any person
	obtaining a copy of this software and associated documentation
	files (the "Software"), to deal in the Software without
	restriction, including without limitation the rights to use,
	copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the
	Software is furnished to do so, subject to the following
	conditions:

	The above copyright notice and this permission notice shall be
	included in all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
	EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
	OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
	NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
	HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
	WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
	FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
	OTHER DEALINGS IN THE SOFTWARE.


