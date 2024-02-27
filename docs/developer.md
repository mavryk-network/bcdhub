## Overview
BCDHub is a set of microservices written in Golang:

* `indexer`  
Loads and decodes operations related to smart contracts and also keeps track of the blockchain and handles protocol updates.
* `API`  
Exposes RESTful JSON API for accessing indexed data (with on-the-fly decoding). Also provides a set of methods for authentication and managing user profiles.

Those microservices are sharing access to databases and communicating via database:

* `PostgreSQL` database for storing compilations and user data.

### Third-party services
BCDHub also depends on several API endpoints exposed by [TzKT](https://github.com/mavryk-network/tzkt) although they are optional:

* List of blocks containing smart contract operations, used for boosting the indexing process (allows to skip blocks with no contract calls)
* Mempool operations
* Contract aliases and other metadata

Those services obviously make sense for public networks only and not used for sandbox or other private environments.

## Versioning
BCD uses `X.Y.Z` version format where:
* `X` changes every 3-5 months along with a big release with a significant addition of functionality  
* `Y` increasing signals about a possibly non-compatible update that requires reindexing (or restoring from snapshot) or syncing with frontend
* `Z` bumped for every stable release candidate or hotfix

### Syncing with frontend
BCD web interface developed at https://github.com/mavryk-network/bcd uses the same version scheme.  
`X.Y.*` versions of backend and frontent MUST BE compatible which means that for every change in API responses `Y` has to be increased.

### Publishing releases
Is essentially tagging commits:
```bash
make release  # forced tag update
```
For stable release:
```bash
git tag X.Y.Z
git push --tags
```

## Docker images
Although you can install and run each part of BCD Hub independently, as system services for instance, the simplest approach is to use dockerized versions orchestrated by _docker-compose_.  

BCDHub docker images are being built on [dockerhub](https://hub.docker.com/u/bakingbad). Tags for stable releases have format `X.Y`.

### Linking with Git tags
Docker tags are essentially produced from Git tags using the following rules:
* `X.Y.*` → `X.Y`

### Building images
```bash
make images  # latest
make stable-images  # requires STABLE_TAG variable in the .env file
```



## Deploy

### Requirements
Make sure you have installed:
* docker
* docker-compose

You will also need several ports to be not busy:
* `14000` API service
* `5432` PostgreSQL
* `8000` Frontend GUI

### Get ready
1. Clone this repo
```bash
git clone https://github.com/mavryk-network/bcdhub.git
cd bcdhub
```

2. Create and fill `.env` file (see _Configuration_)
```bash
your-text-editor .env
```

### Environments
There are several predefined configurations serving different purposes.

#### Production `better-call.dev`
* Stable docker images `X.Y`
* `/configs/production.yml` file is used internally
* Requires `STABLE_TAG` environment set
* Deployed via `make stable`

#### Development `localhost`
* `/configs/development.yml` file is used
* You can spawn local instances of databases or _ssh_ to staging host with port forwarding
* Run services `make {service}` (where service is one of `api` `indexer`)

#### Sandbox `bcdbox`
* `/configs/sandbox.yml` file is used
* Start via `COMPOSE_PROJECT_NAME=bcd-box docker-compose -f docker-compose.sandbox.yml up -d --build`
* Stop via `COMPOSE_PROJECT_NAME=bcd-box docker-compose -f docker-compose.sandbox.yml down`


## Running

### Startup
It takes around 20-30 seconds to initialize all services, API endpoints might return errors until then.  
**NOTE** that if you specified local RPC node that's not running, BCDHub will wait for it indefinitely.

## Snapshots
Full indexing process requires about 2 hours, however there are cases when you cannot afford that

**NOTE:** currently we don't provide public snapshots.
Alternatively, contact us for granting access

### Get ready
* Make sure you have snapshot settings in your `.env` file

### Make snapshot

#### 1. Initialize credentials
```
make s3-creds
```
No further actions required

#### 2. Create local repository (if not exists)
```
make s3-repo
```
Follow the instruction: you can choose an arbitrary name for your repo.

#### 3. Create snapshot
```
make s3-snapshot
```
Select an existing repository to store your snapshot.

### Restore snapshot

#### 1. Initialize creds and repo
Follow steps 1 and 2 from the _make snapshot_ instruction.

#### 2. Apply snapshot
```
make s3-restore
```
Select the latest (by date) snapshot from the list. It's taking a while, don't worry about the seeming freeze.

## Version upgrade
This is mostly for production environment, for all others a simple "start from the scratch" would work.

### Soft update
E.g. applying hotfixes. No breaking changes in the database schema.

#### 1. Build stable images
Make sure you are on master branch
```
git pull
make stable-images
make stable
```

#### 1'. Pull stable images
```
make stable-pull
```

#### 2. Deploy
```
make stable
```

### Data migration
E.g. new field added to one of the models. You'd need to write a migration script to update existing data.

#### 1. Pull migration script
```
git pull
```

#### 2. Run migration
```
make migration
```
Select your script.


### Upgrade from snapshot
In case you need to reindex from scratch you can set up a secondary BCDHub instance, fill the index, make a snapshot, and then apply it to the production instance.

#### 0. Make a snapshot
Typically you'd use staging for that.

#### 1. Stop BCDHub and clear indexed data
```
make upgrade
```

#### 2. Restore snapshot
```
make s3-restore
```
Select the snapshot you made.

#### 3. Run the rest of the services
```
make stable
```
