# SPA Clock

SPA Clock is a small Single Page Application (SPA) clock that runs
well on a Raspberry Pi touch screen. Both the 2.8" and 6" Touch
Panels. 

>It will actually run on pretty much any browser, but it is more fun
>when it consumes a touch screen on a tiny computer! :)

## Build and Use spaclock

This is a single go file, with additional web files and hence can be
run in your favorite way go way.  To build and run just do:

> go build && ./spaclock 

To get help running the program use the -help / -h switch:

> ./spaclock -h

Now point your browser at http://spaclock and that is it.

## Building SPA Clock

### Cut-n-Paste Internet Clock

This first version of SPA clock is extrodinarily simple, I literally
copy and pasted the following snippets of code to create a very
limited but useful Internet Clock:

1. [Gorilla Mux SPA](https://github.com/gorilla/mux) the example code
  in the _Serving Single Page Applications_ section.

2. [Bootstrap index.html
  template](https://getbootstrap.com/docs/4.4/getting-started/introduction/#starter-template)
  Look for section Starter Template, I cut and pasted the index.

3. [JavaScript Clock from
w3school](https://www.w3schools.com/js/tryit.asp?filename=tryjs_timing_clock),
I combined the code between the script tags with index.html 

### Self Resetting Clock with NTP!

This simple clock connects to the internet getting its time from a
public NTP server, so no need to ever reset the clock again. 

### Style Matters

In addition to the copy and pasted code above, I added a little bit
(very tiny bit of not good) styling by changing the background color
and choosing a [google font](http://google.com/fonts). Just to perk
the imagine for styling possibilities.

## Moving Forward

It does not do much more than that at present, however I can imagine a
couple cool and very useful features that I would like to add to our
_internet clock_, like display different timezones, alarms and display
the weather, for example.

Now don't let me get too far ahead of myself, I have quite a bit of
work to do laying a couple important pieces of plumbing, we'll get
into the relavent issues shortly, as they come up!
