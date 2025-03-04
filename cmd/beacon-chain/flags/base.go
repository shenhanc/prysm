// Package flags defines beacon-node specific runtime flags for
// setting important values such as ports, eth1 endpoints, and more.
package flags

import (
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/urfave/cli/v2"
)

var (
	// HTTPWeb3ProviderFlag provides an HTTP access endpoint to an ETH 1.0 RPC.
	HTTPWeb3ProviderFlag = &cli.StringFlag{
		Name:  "http-web3provider",
		Usage: "A mainchain web3 provider string http endpoint. Can contain auth header as well in the format --http-web3provider=\"https://goerli.infura.io/v3/xxxx,Basic xxx\" for project secret (base64 encoded) and --http-web3provider=\"https://goerli.infura.io/v3/xxxx,Bearer xxx\" for jwt use",
		Value: "",
	}
	// FallbackWeb3ProviderFlag provides a fallback endpoint to an ETH 1.0 RPC.
	FallbackWeb3ProviderFlag = &cli.StringSliceFlag{
		Name:  "fallback-web3provider",
		Usage: "A mainchain web3 provider string http endpoint. This is our fallback web3 provider, this flag may be used multiple times.",
	}
	// DepositContractFlag defines a flag for the deposit contract address.
	DepositContractFlag = &cli.StringFlag{
		Name:  "deposit-contract",
		Usage: "Deposit contract address. Beacon chain node will listen logs coming from the deposit contract to determine when validator is eligible to participate.",
		Value: params.BeaconConfig().DepositContractAddress,
	}
	// RPCHost defines the host on which the RPC server should listen.
	RPCHost = &cli.StringFlag{
		Name:  "rpc-host",
		Usage: "Host on which the RPC server should listen",
		Value: "127.0.0.1",
	}
	// RPCPort defines a beacon node RPC port to open.
	RPCPort = &cli.IntFlag{
		Name:  "rpc-port",
		Usage: "RPC port exposed by a beacon node",
		Value: 4000,
	}
	// MonitoringPortFlag defines the http port used to serve prometheus metrics.
	MonitoringPortFlag = &cli.IntFlag{
		Name:  "monitoring-port",
		Usage: "Port used to listening and respond metrics for prometheus.",
		Value: 8080,
	}
	// CertFlag defines a flag for the node's TLS certificate.
	CertFlag = &cli.StringFlag{
		Name:  "tls-cert",
		Usage: "Certificate for secure gRPC. Pass this and the tls-key flag in order to use gRPC securely.",
	}
	// KeyFlag defines a flag for the node's TLS key.
	KeyFlag = &cli.StringFlag{
		Name:  "tls-key",
		Usage: "Key for secure gRPC. Pass this and the tls-cert flag in order to use gRPC securely.",
	}
	// DisableGRPCGateway for JSON-HTTP requests to the beacon node.
	DisableGRPCGateway = &cli.BoolFlag{
		Name:  "disable-grpc-gateway",
		Usage: "Disable the gRPC gateway for JSON-HTTP requests",
	}
	// GRPCGatewayHost specifies a gRPC gateway host for Prysm.
	GRPCGatewayHost = &cli.StringFlag{
		Name:  "grpc-gateway-host",
		Usage: "The host on which the gateway server runs on",
		Value: "127.0.0.1",
	}
	// GRPCGatewayPort specifies a gRPC gateway port for Prysm.
	GRPCGatewayPort = &cli.IntFlag{
		Name:  "grpc-gateway-port",
		Usage: "The port on which the gateway server runs on",
		Value: 3500,
	}
	// EthApiPort specifies the port which runs the official Ethereum REST API.
	// Serves JSON values conforming to the specification: https://ethereum.github.io/eth2.0-APIs/
	EthApiPort = &cli.IntFlag{
		Name:  "eth-api-port",
		Usage: "The port which exposes a REST API conforming to the official Ethereum API specification.",
		Value: 3501,
	}
	// GPRCGatewayCorsDomain serves preflight requests when serving gRPC JSON gateway.
	GPRCGatewayCorsDomain = &cli.StringFlag{
		Name: "grpc-gateway-corsdomain",
		Usage: "Comma separated list of domains from which to accept cross origin requests " +
			"(browser enforced). This flag has no effect if not used with --grpc-gateway-port.",
		Value: "http://localhost:4200,http://localhost:7500,http://127.0.0.1:4200,http://127.0.0.1:7500,http://0.0.0.0:4200,http://0.0.0.0:7500",
	}
	// MinSyncPeers specifies the required number of successful peer handshakes in order
	// to start syncing with external peers.
	MinSyncPeers = &cli.IntFlag{
		Name:  "min-sync-peers",
		Usage: "The required number of valid peers to connect with before syncing.",
		Value: 3,
	}
	// ContractDeploymentBlock is the block in which the eth1 deposit contract was deployed.
	ContractDeploymentBlock = &cli.IntFlag{
		Name:  "contract-deployment-block",
		Usage: "The eth1 block in which the deposit contract was deployed.",
		Value: 11184524,
	}
	// SetGCPercent is the percentage of current live allocations at which the garbage collector is to run.
	SetGCPercent = &cli.IntFlag{
		Name:  "gc-percent",
		Usage: "The percentage of freshly allocated data to live data on which the gc will be run again.",
		Value: 100,
	}
	// HeadSync starts the beacon node from the previously saved head state and syncs from there.
	HeadSync = &cli.BoolFlag{
		Name:  "head-sync",
		Usage: "Starts the beacon node with the previously saved head state instead of finalized state.",
	}
	// SlotsPerArchivedPoint specifies the number of slots between the archived points, to save beacon state in the cold
	// section of DB.
	SlotsPerArchivedPoint = &cli.IntFlag{
		Name:  "slots-per-archive-point",
		Usage: "The slot durations of when an archived state gets saved in the DB.",
		Value: 2048,
	}
	// DisableDiscv5 disables running discv5.
	DisableDiscv5 = &cli.BoolFlag{
		Name:  "disable-discv5",
		Usage: "Does not run the discoveryV5 dht.",
	}
	// BlockBatchLimit specifies the requested block batch size.
	BlockBatchLimit = &cli.IntFlag{
		Name:  "block-batch-limit",
		Usage: "The amount of blocks the local peer is bounded to request and respond to in a batch.",
		Value: 64,
	}
	// BlockBatchLimitBurstFactor specifies the factor by which block batch size may increase.
	BlockBatchLimitBurstFactor = &cli.IntFlag{
		Name:  "block-batch-limit-burst-factor",
		Usage: "The factor by which block batch limit may increase on burst.",
		Value: 10,
	}
	// DisableSync disables a node from syncing at start-up. Instead the node enters regular sync
	// immediately.
	DisableSync = &cli.BoolFlag{
		Name:  "disable-sync",
		Usage: "Starts the beacon node without entering initial sync and instead exits to regular sync immediately.",
	}
	// EnableDebugRPCEndpoints as /v1/beacon/state.
	EnableDebugRPCEndpoints = &cli.BoolFlag{
		Name:  "enable-debug-rpc-endpoints",
		Usage: "Enables the debug rpc service, containing utility endpoints such as /prysm/v1alpha1/beacon/state.",
	}
	// SubscribeToAllSubnets defines a flag to specify whether to subscribe to all possible attestation subnets or not.
	SubscribeToAllSubnets = &cli.BoolFlag{
		Name:  "subscribe-all-subnets",
		Usage: "Subscribe to all possible attestation subnets.",
	}
	// HistoricalSlasherNode is a set of beacon node flags required for performing historical detection with a slasher.
	HistoricalSlasherNode = &cli.BoolFlag{
		Name:  "historical-slasher-node",
		Usage: "Enables required flags for serving historical data to a slasher client. Results in additional storage usage",
	}
	// ChainID defines a flag to set the chain id. If none is set, it derives this value from NetworkConfig
	ChainID = &cli.Uint64Flag{
		Name:  "chain-id",
		Usage: "Sets the chain id of the beacon chain",
	}
	// NetworkID defines a flag to set the network id. If none is set, it derives this value from NetworkConfig
	NetworkID = &cli.Uint64Flag{
		Name:  "network-id",
		Usage: "Sets the network id of the beacon chain.",
	}
	// WeakSubjectivityCheckpt defines the weak subjectivity checkpoint the node must sync through to defend against long range attacks.
	WeakSubjectivityCheckpt = &cli.StringFlag{
		Name: "weak-subjectivity-checkpoint",
		Usage: "Input in `block_root:epoch_number` format. This guarantees that syncing leads to the given Weak Subjectivity Checkpoint along the canonical chain. " +
			"If such a sync is not possible, the node will treat it a critical and irrecoverable failure",
		Value: "",
	}
	// Eth1HeaderReqLimit defines a flag to set the maximum number of headers that a deposit log query can fetch. If none is set, 1000 will be the limit.
	Eth1HeaderReqLimit = &cli.Uint64Flag{
		Name:  "eth1-header-req-limit",
		Usage: "Sets the maximum number of headers that a deposit log query can fetch.",
		Value: uint64(1000),
	}
	// GenesisStatePath defines a flag to start the beacon chain from a give genesis state file.
	GenesisStatePath = &cli.StringFlag{
		Name: "genesis-state",
		Usage: "Load a genesis state from ssz file. Testnet genesis files can be found in the " +
			"eth2-clients/eth2-testnets repository on github.",
	}
)
