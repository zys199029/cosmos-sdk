# Crypto

Originally in github.com/tendermint/go-crypto, has been merged into
github.com/tendermint/tendermint/crypto. 

## Hash Function

Call it
[`TMHASH`](https://github.com/tendermint/tendermint/blob/develop/crypto/tmhash/hash.go). 
We recently changed it from RIPEMD160 to SHA256[:20], 
though there's concern that 80-bits of SHA2 is unsafe given bitcoin miners.
We use TMHASH for merkle trees and for computing addresses from pubkeys.

Relevant issues:
- Cost to attack SHA2 [#1990](https://github.com/tendermint/tendermint/issues/1990)
- Switch to SHA2 [iavl/#38](https://github.com/tendermint/iavl/issues/38)

## Merkle Tree

Two tree types: 
[Simple](https://github.com/tendermint/tendermint/tree/develop/crypto/merkle) and 
[IAVL](https://github.com/tendermint/iavl).

Simple is just a simple static Merkle tree used for hashing things
in blocks like transactions, validator sets, and headers. 

IAVL is the dynamic immutable AVL+ tree for holding application state.
It's not actually used directly by Tendermint, but by an app that needs it.

Both use TMHASH.

Relevant Issues
- Prefix internal nodes in the Simple tree, perhaps adopt RFC 6962 [#1892](https://github.com/tendermint/tendermint/issues/1892)

## Crypto Deps

We use independent forks of bcrypt and ed25519 - we should consolidate and just fork
golang/x/crypto. Reason for the forks is typically the ability to inject another source 
of randomness (for those that don't trust OS entropy...)

See [#1959](https://github.com/tendermint/tendermint/issues/1959)

## Authenticated Encryption

We use Station-to-Station protocol for authenticated encryption handshake between peers using
NACL.SecretBox, which is xsalsa20poly1305. This is being changed to use hkdf-chacha20poly1305 (see [#1563](https://github.com/tendermint/tendermint/pull/1563)).
Note some matters appear yet unresolved in how to pull material from the hkdf. See issues:

- Use chacha instead of salsa [#1124](https://github.com/tendermint/tendermint/issues/1124)
- Use hkdf instead of SecretBox API with kludgy nonce [#1165](https://github.com/tendermint/tendermint/issues/1165)

Note we use this for connecting to peers and for the connection between Tendermint Core and a validator signing process.

## Keys

We use NACL.SecretBox for symetric encryption of user keys. See
[symmetric.go](https://github.com/tendermint/tendermint/blob/develop/crypto/xsalsa20Symmetric/symmetric.go).
Note there's the suggestion to switch from NACL.SecretBox (xsalsa20poly1305) to hkdf-chacha20poly1305 
[#1955](https://github.com/tendermint/tendermint/issues/1955) for the key encryption too (ie. to match the p2p auth enc).

The SDK houses a client pkg that does all the key management, and uses the encryption from tendermint/crypto.
It also supports the Ledger Nano S.


We are committing to supporting Secp256k1 for user keys, and Ed25519 for validator keys for now. 
For users, see the
[cosmos-sdk/client](https://github.com/cosmos/cosmos-sdk/tree/develop/client) 
and [cosmos-sdk/crypto](https://github.com/cosmos/cosmos-sdk/tree/develop/crypto) packages.
For validators, see the [tendermint/privval](https://github.com/tendermint/tendermint/tree/develop/privval) package. 
Note we intend to use an external process for validators called the [KMS](https://github.com/tendermint/kms) that's being written in Rust. It will support multiple HSM backends.

# Tendermint

Previous reviews/tests:

- Jepsen serializability tests. [See report](https://jepsen.io/analyses/tendermint-0-10-2)
    - Summer, 2017
    - No safety or liveness violations found in consensus under adverse network conditions 
        and byzantine validators (ie. simple duplicate validators)
    - Though found we weren't fsyncing the consensus WAL. Fixed
- Independent Auditor 1
    - Sept, 2017
    - Mostly focused on crypto and serialization layer
- Internal Auditor 2
    - Sept 2017 - March 2018
    - Focused on fuzzing, see [the repo](https://github.com/tendermint/fuzz).
        - Primarily on the underlying data structures and API
        - Need to now focus on fuzzing the reactor messages
- Zarko Milosevic (Internal Researcher)
    - Nov 2017
    - BFT analysis - found liveness bug
    - Fixed with modifications to protocol. [See recent paper](https://arxiv.org/abs/1807.04938)
        - Formal proofs of safety and liveness
    - Still need to make the implementation 100% match this spec
- Formal Modelling (External Researchers)
    - Early 2018. 
    - Proved safety, but found same liveness bug as Zarko
    - Working to improve the tool to check the fixed version
- Independent Auditor 2
    - March 2018. Reviewed the crypto, libs, and some of Tendermint Core
    - Mostly minor issues found
- Independent Auditor 3
    - May 2018. Reviewed the primary interfaces/sockets/endpoints
    - Some serious p2p issues found, still need to address (!)
- HackerOne Bug Bounty
    - Started in May 2018
    - Just found a serious bug in one of the RPC endpoints - still need to address (!)
- Ongoing Testnets
    - Testnets have been getting larger and seriously stressing the Tendermint components
    - gaia-5001 had bug in evidence reactor that led to very high peer churn
        - we kept making blocks with ~20-30 validators through it all
    - gaia-6002 went for ~1 month without real issue. started with ~30 validators, went up to ~50
    - gaia-7001 launched with 107 validators in genesis. Currently has 130.
        - just crashed with a panic in the state machine

Security concerns:
- reactor messages that could cause panic, OOM, amplification attack, etc.
- rpc endpoints that could cause panic, OOM, amplification attack, etc.
- p2p eclipse attacks and better tracking of bad peers
- Also see [issues labelled with security tag](https://github.com/tendermint/tendermint/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3Asecurity+)
    - probably others should be labelled "security" as well

Other general issues:
- P2P configuration is confusing and needs a refactor
- AddrBook could be improved - make use of chain-id, DNS, bad peerness
- Large inefficiencies in the mempool [#1798](https://github.com/tendermint/tendermint/issues/1798)
- Some inefficiencies in the blockchain reactor (used to fast-sync the chain)

Open TODO
- Better observability 
- BFT Time [#2013](https://github.com/tendermint/tendermint/issues/2013)
- NextValSet - validator set changes from block H take place at H+2 [#1815](https://github.com/tendermint/tendermint/pull/1815)
- Minor ABCI and header changes. [See all ABCI
  tags](https://github.com/tendermint/tendermint/issues?q=is%3Aissue+is%3Aopen+label%3Aabci+milestone%3Alaunch)
- Add protocol version [#1983](https://github.com/tendermint/tendermint/pull/1983)

# Amino

[Amino](https://github.com/tendermint/go-amino) is our serialization library. It's proto3 without `oneof` but with a concept of `interfaces` 
which use prefix bytes that are the hash of the registered type name.

Previous versions were reviewed (before it was proto3 compatible, and also back when it was called go-wire).
The latest proto3 compatible version needs review.

# SDK

[Cosmos-SDK](https://github.com/cosmos/cosmos-sdk) is the Go framework for building apps.
For a quick introduction to the SDK, see the [intro guide](https://cosmos.network/docs/sdk/core/intro.html).

The following comprise the rough core of the SDK and are ready for review:

- `store` pkg, which contains the MultiStore and the logic for caching and commiting stores.
- `types` pkg, which contains various concrete types like big integers, rational numbers, coins, 
    and various interface/function types like handlers and messages.
- `baseapp` pkg, which abstracts over the ABCI and connects message and their handlers with the stores.
- `x/auth` pkg, which defines the default tx type and how it is authenticated
- `x/bank` pkg, which defines the basic coin transfer mechanism

The remaining modules and client pkgs are still under more active development and internal review (see Gaia, below).

Note Gas counting has been threaded through the SDK but is not well specified yet and still needs to be measured/tweaked.

# Gaia

This is the binary for the Cosmos Hub.
It includes the `x/stake`, `x/gov`, and other `x/` modules found in the SDK, which contain
all the logic for the Cosmos Hub staking, slashing, governance, etc.
While x/stake and x/gov are roughly feature complete, they are still being internally reviewed and improved (we continue to catch them panicing on testnets).

That said, the [specs](https://github.com/cosmos/cosmos-sdk/tree/develop/docs/spec) 
have a significant amount of detail that are worthy of review.

Note the first few specs have not been written (Store, Auth, Bank), but Staking, Governance, 
and Slashing are relatively complete, while Provisioning (ie. fee distribution and inflationary rewards) 
and IBC have significant detail but are still being iterated on.

We still haven't implemented fee withdrawals or the special logic for eg. AiB vesting its Atoms
