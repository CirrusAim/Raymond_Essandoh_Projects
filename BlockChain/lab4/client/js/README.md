# Native javascript client

This example was built using only javascript and [webpack](https://webpack.js.org/), and it connects with the blockchain through the [web3](https://github.com/ethereum/web3.js/) API.

## Installing the necessary dependencies
```
npm install
```

## Build the example
```
npm run build
```

The command above will generate the `dist` directory with your application. We use webpack to bundle all the dependencies and generate only one javascript (i.e. `app.js`) that is used in the `index.html`.

## Running the example

You must have an Ethereum node running to connect your application.
By default the example application attempts to connect to `http://127.0.0.1:7545` with the network ID `5777`, and you can start the `ganache-cli` by running the npm command from the root directory of [lab 4](../../README.md).
You can also use an official Ethereum Testnet like [Ropsten](https://ropsten.etherscan.io/) if you want, but you will be required to fund some accounts using a [faucet](https://faucet.ropsten.be/).

1. Running the development blockchain environment from the `lab4`

+ Open on terminal and go to the `lab4` directory and run:
```
npm run ganache-cli
```

+ Deploy the contracts in the running blockchain instance:

```
npm run migrate:ganache
```

2. Serving the web front-end on http://localhost:8080

Open another terminal and from the `lab4/client/js` and run:
```
npm run dev
```

Or, to build and then start the server:

```
npm run build
npm run serve
```

3. Open your web browser at: http://localhost:8080
