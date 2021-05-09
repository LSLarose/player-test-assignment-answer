module player-updater

go 1.16

//list all local modules used
replace github.com/LSLarose/httpPlayerServer => ./httpPlayerServer
replace github.com/LSLarose/greetings => ./greetings

require (
	github.com/LSLarose/httpPlayerServer v0.0.0-00010101000000-000000000000
	rsc.io/quote/v3 v3.1.0 // indirect
)
