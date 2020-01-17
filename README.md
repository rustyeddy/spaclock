# SPA Clock

SPA Clock is a small Single Page Application (SPA) iClock that runs
well on a Raspberry Pi touch screen. Both the 2.8" and 6" Touch
Panels. It will actually run on pretty much any browser, but it is
more fun when it consumes the full touch screen on a tiny computer! :) 

This is the first hack at creating something useful, we'll call it a
_Minimum Viable Clock_. 

## SPA Clock Features

OK, spaClock is an Internet Clock meaning it is a clock connected to
the Internet, giving the clock access to a tremendous amount of
resources allow it to provide some really cool features like:

- Self setting with the Internet and NTP!
- Message section to send yourself all sorts of messages
- Websockets for realtime updates, without refreshing!

This simple clock connects to the internet getting its time from a
public NTP server, so no need to reset the clock, ever again, even
across changes for Daylight Savings! 

Expect more cool stuff to come along.  Hmm. Need to make a plugin
system??? :)

## Building The spaClock

This is a single go file, with additional web files and hence can be
run in your favorite way go way.  To build and run just do:

> go build && ./spaclock 

To get help running the program use the -help / -h switch:

> ./spaclock -h

Now point your browser at http://host-or-ip:8000/ and that is it! 


### SPA Clock Design

SPAClock is a Single Page App (SPA) with the _backend written in go_.
Since SPAClock is a SPA (uuuhhh) it derives the UI from HTML5/CSS and
JavaScript. The _backend_ is a self contained server written in
[go](http://golang.org) ready to serve up static HTML, CSS &
JavaScript, it also has a REST API and Sports a Websocket for real
time communication.

## Developing SPAClock

### Cut-n-Paste Internet Clock

This first version of SPA clock was extrodinarily simple, I literally
copy and pasted the following snippets from their respective internet
sources. 

The code has since become more complex for reasons we'll get into
later, but the skeleton actually provided a very nice starting point! 


1. [Gorilla Mux SPA](https://github.com/gorilla/mux) the example code
  in the _Serving Single Page Applications_ section.

2. [Bootstrap index.html
  template](https://getbootstrap.com/docs/4.4/getting-started/introduction/#starter-template)
  Look for section Starter Template, I cut and pasted the index.

3. [JavaScript Clock from
w3school](https://www.w3schools.com/js/tryit.asp?filename=tryjs_timing_clock),

The code snippets above were combined into _main.go_, _index.html_ as
well as created the _clock.js_ file.  I added bootstrap to give app a
littl style and layout. 

I am not a designer, so I decided to keep the clock simple:

### Style Matters

In addition to the copy and pasted code above, I added a little bit
(very tiny bit of not good) styling by changing the background color
and choosing a [google font](http://google.com/fonts). Just to perk
the imagine for styling possibilities.

## Moving Forward

- Alarms
- Weather (weather.org)
- MQTT to display sensor and other daya
- Plugin system
