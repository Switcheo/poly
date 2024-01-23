module github.com/polynetwork/poly

go 1.14

require (
	github.com/Zilliqa/gozilliqa-sdk v1.2.1-0.20210329093354-1b8e0a7a2e25
	github.com/btcsuite/btcd v0.23.0
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/cosmos/cosmos-sdk v0.47.0-rc2.0.20230220103612-f094a0c33410
	github.com/ethereum/go-ethereum v1.9.15
	github.com/gcash/bchd v0.16.5
	github.com/gcash/bchutil v0.0.0-20200506001747-c2894cd54b33
	github.com/gorilla/websocket v1.5.0
	github.com/gosuri/uiprogress v0.0.1
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/howeyc/gopass v0.0.0-20190910152052-7cb4b85ec19c
	github.com/itchyny/base58-go v0.1.0
	github.com/joeqian10/neo-gogogo v1.1.0
	github.com/joeqian10/neo3-gogogo v0.3.4
	github.com/ontio/ontology v1.11.1-0.20200812075204-26cf1fa5dd47
	github.com/ontio/ontology-crypto v1.0.9
	github.com/ontio/ontology-eventbus v0.9.1
	github.com/pborman/uuid v1.2.0
	github.com/polynetwork/poly-io-test v0.0.0-20200819093740-8cf514b07750
	github.com/stretchr/testify v1.8.4
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d
	github.com/tendermint/tendermint v0.34.9
	github.com/urfave/cli v1.22.4
	github.com/valyala/bytebufferpool v1.0.0
	golang.org/x/crypto v0.11.0
	golang.org/x/net v0.12.0
	gotest.tools v2.2.0+incompatible
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/cosmos/cosmos-sdk => github.com/Switcheo/cosmos-sdk v0.47.5-0.20240119065259-675e01adc46f

replace github.com/btcsuite/btcd => github.com/btcsuite/btcd v0.22.2
