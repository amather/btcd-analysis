module github.com/amather/btcd-analysis

go 1.17

replace github.com/btcsuite/btcd => ../btcd
replace github.com/amather/shanalyzer => ../shanalyzer

require (
	github.com/btcsuite/btcd v0.20.1-beta
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/amather/shanalyzer v0.0.0-ignore
)

require (
	github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
	github.com/btcsuite/go-socks v0.0.0-20170105172521-4720035b7bfd // indirect
	github.com/btcsuite/websocket v0.0.0-20150119174127-31079b680792 // indirect
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
)
