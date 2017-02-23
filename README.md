# GoBackLight
Golang based backlight tool

This tool is installed with setuid root to allow me to update my laptop backlight
without having to run sudo.

Calling syntax:

````backlight [-v] [+-]<int>[%]````

-v will enable verbose mode. If specified, this needs to be the first argument.

The value argument (second if -v appears, first otherwise) accepts an initial +
or - to indicate that you want a change relative to the current value, and accepts
a % to indicate your value is specified in a percentage.
