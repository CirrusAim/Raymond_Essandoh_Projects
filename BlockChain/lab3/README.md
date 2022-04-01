# Lab 3: Introduction to Smart Contracts

| Lab 3:           | Smart Contracts              |
| ---------------- | ---------------------------- |
| Subject:         | DAT650 Blockchain Technology |
| Deadline:        | 02. NOV                      |
| Expected effort: | 2 weeks                      |
| Grading:         | Pass/fail                    |

## Table of Contents
- [Lab 3: Introduction to Smart Contracts](#lab-4-introduction-to-smart-contracts)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Tools](#tools)
    - [Install and Running](#install-and-running)
  - [Lab Approval](#lab-approval)

## Introduction

In this lab you need to write two smart contracts.
The description of each part of the assignment can be found under the respective directories, [Wallet contract](wallet/README.md) and [Betting contract](betting/README.md).
All tests should pass, amd you can add as many functions as you need.

## Tools

For this lab you will need to have [Truffle Suite](https://www.trufflesuite.com/docs/truffle/overview) installed in your machine.
It is also desirable that you have [Ganache](https://www.trufflesuite.com/docs/ganache/overview) installed or any other develop blockchain configured, with at least two accounts, to perform correctly the tests.

Both assignments have a `package.json` file with the dependencies and scripts for easy development, including truffle and ganache-cli, that can be installed locally using the `npm`. Familiarize yourself with it can save you some time while developing, you can start [here](https://nodejs.dev/learn/the-package-json-guide/#scripts).

* Note that the commands shown below and specified in the scripts section of the `package.json` file are optional. If you have truffle installed globally in your system you can use it instead, by running direct the commands specified in the scripts.

### Install and Running

To install the necessary dependencies to run and test each assignment, enter in the correspondent directory and run the `npm install` command. Like in the example below for the wallet project:

```
$ cd wallet
$ npm install
```

After the installation you can compile and run the tests as following:
```
$ npm run compile
$ npm run migrate:ganache
$ npm run test:ganache
```

If you get the following error:

```
Could not connect to your Ethereum client with the following parameters:
    - host       > 127.0.0.1
    - port       > 8545
    - network_id > *
Please check that your Ethereum client:
    - is running
    - is accepting RPC connections (i.e., "--rpc" option is used in geth)
    - is accessible over the network
    - is properly configured in your Truffle configuration file (truffle-config.js)
```

It means that you need to have running a blockchain instance in another terminal.
There are many options to perform that, and you can find more information [here](https://www.trufflesuite.com/docs/truffle/reference/choosing-an-ethereum-client).

For the purpose of this lab, we will be using the [ganache](https://www.trufflesuite.com/docs/ganache/overview).
You can use the [ganache GUI](https://github.com/trufflesuite/ganache) or the [ganache-cli](https://github.com/trufflesuite/ganache-cli/blob/master/README.md) command-line tool, both with same setup.
More information about the ganache configuration can be found [here](https://www.trufflesuite.com/docs/ganache/truffle-projects/linking-a-truffle-project)

A pre-configured ganache-cli is available in the _scripts_ section of the `package.json`. You can use it by running:
```
$ npm run ganache-cli
```

## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours and present your solution. This lets you present the thought process behind your solution, and allows us to provide feedback on your solution then and there.
When you are ready to show your solution, reach out to a member of the teaching staff. It is expected that you can explain your code and show how it works. You may show your solution on a lab workstation or your own computer.

You should demonstrate that your implementation fulfills the previously listed specification of each assignments part.
The task will be verified by a member of the teaching staff during lab hours.
