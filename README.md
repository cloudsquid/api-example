# Cloudsquid API example usage
Example usage of the cloudsquid api

## First steps

Follow the steps in our [documentation](https://docs.cloudsquid.io/docs/quickstart) to create an account and set up your first extraction agent. You don't need to add any files or start runs here

## Golang example
The relevant API logic is in `go/main.go`

Create a copy of `go/.env.example` name it `.env` and add the relevant environmental variables in the new `go/.env`
- `CLOUDSQUID_API_KEY`: Your API key, which you can create in the cloudsquid dashboard settings
- `CLOUDSQUID_API_ENDPOINT`: The URL of the cloudsquid API. This is `https://api.cloudsquid.io/api`
- `CLOUDSQUID_AGENT_ID`: The ID of the extraction agent you created in the cloudsquid dashboard. You can find it on top of the agent page you intend to use. You can also find this as the identifier of the URL slug when you are on the agent page. 

After this run `cd go` and `go run . -f <path to your file>` to run the example. The file will be uploaded to cloudsquid and the extraction agent will start running. Results will be logged to your terminal. You can see the results in the cloudsquid ui!

## JS example
The relevant API logic is in `js/index.js`

First navigate to the js folder `cd js`

1. Run `npm install` to install the dependencies
2. Create a copy of `js/.env.example`, rename to `.env` and add the relevant environmental variables in the new  `js/.env` (see golang example)
3. Run `npm start <path_to_your_file>` to start the example.

