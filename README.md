# IWS_GO
## Preperation
For this repository to run, you need to install Go.
You also need to setup a MongoDB database and provide a valid connection string.
For that, rename "connection_CHANGE.json" to "connection.json" and fill in your connection string.

## Installation
To install all dependencies after cloning this repository, you can run the command:
`go mod tidy`

To install the app from this repository run:
`go install`

## Execution
An executable should be placed in the /bin directory of the go installation.
Navigate there (or add that path to your PATH variable) and execute the app:
`IWS_GO`