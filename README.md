# Puppet

A visual command-line interface for managing, creating, and analyzing SummerCash networks.

## Installation

```zsh
go install github.com/SummerCash/puppet
```

## Usage

### Creating a New SummerCash Network

```zsh
puppet create
```

### Searching for Data In the SummerCash Blockmesh

```zsh
puppet search SEARCH_TERM BLOCKCHAIN_TO_SEARCH_IN SECOND_BLOCKCHAIN_TO_SEARCH_IN
```

Note: By only providing a SEARCH_TERM, puppet will search the entire blockmesh, rather than a single chain or group of chains.
