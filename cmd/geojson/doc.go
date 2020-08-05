/*
Geojson provides utilities for OpenStreetMap

USAGE:
   geojson [global options] command [command options] [arguments...]

COMMANDS:
   serve    serve the web server
   subarea  list all sub-areas of an OpenStreetMap object
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --out value, -o value  specify the directory of outputs (default: "./geo")
   --verbose              enable verbose logging with DEBUG level (default: false)
   --help, -h             show help (default: false)
   --version, -v          print the version (default: false)


SUBAREA

geojson subarea - list all sub-areas of an OpenStreetMap object

USAGE:
   geojson subarea [command options] [arguments...]

OPTIONS:
   --raw, -r        leave tags in unfornalized form (UNF) (default: false)
   --separated, -s  leave sub-areas unmerged (default: false)
   --help, -h       show help (default: false)

SERVE

geojson serve - serve the web server

USAGE:
   geojson serve [command options] [arguments...]

OPTIONS:
   --address value, --addr value  set the serving address (default: "127.0.0.1:8181")
   --origin value                 set the CORS origin (default: "*")
   --rate value                   set request-per-second for rate-limiting (default: 10)
   --rate-burst value             set burst size (concurrent requests) for rate-limiting (default: 5)
   --rate-ttl value               set the rate limit TTL for inactive sessions (default: "2m")
   --prefix value                 set static fs handler base path (default: "/static")
   --help, -h                     show help (default: false)
*/
package main
