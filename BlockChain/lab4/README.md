# Lab 4: Developing a DApp

| Lab 4:           | Distributed Application      |
| ---------------- | ---------------------------- |
| Subject:         | DAT650 Blockchain Technology |
| Deadline:        | 10. NOV                      |
| Expected effort: | 1 week                       |
| Grading:         | Pass/fail                    |

## Table of Contents

- [Lab 4: Developing a DApp](#lab-5-developing-a-dapp)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Building the contracts (backend)](#building-the-contracts-backend)
  - [Client application](#client-application)
  - [Lab Approval](#lab-approval)

## Introduction

In this lab you will develop a client application for the [Betting contract](../lab3/betting/README.md) created in the [previous lab](../lab3/README.md).

## Building the contracts (backend)

In the lab4 folder, run the following commands:

1. Installing the necessary dependencies.

```
npm install
```

2. Running the development blockchain environment.

```
npm run ganache-cli
```

3. Compiling and deploying the contracts.

```
npm run compile
npm run migrate:ganache
```

Ensure that you copy your own implementation of the [MyWallet](https://github.com/DAT650-2021/assignments/blob/main/lab4/contracts/MyWallet.sol) contract to `contracts` folder before compile and deploy it. Follow the instructions of each client example implementation.

## Client application

A client implementation for the [Wallet contract](../lab3/wallet/README.md) is given as an example in different languages: in javascript under the directory [client/js](client/js/README.md) and in go under the directory [client/go](client/go/README.md).

The given DApp is only an example, and you are not required to follow the exact setup for your contract.
You can use any framework and language you are conformable with, but you should be able to demonstrate all the functionalities of the Betting contract in your DApp.
If you want to reuse the example setup, you will need to copy your `Betting` contract to the `lab4/contracts` folder and adjust the migrations in `lab4/migrations` folder to deploy your `Betting` contract instead of `MyWallet` (please take a look at the file [2_deploy_contract.js](./migrations/2_deploy_contract.js) and create your own javascript migration file for your contract).

You will also need to create your own client under the `lab4/client` folder and setup it accordingly to use the compiled contract code.
The compiled contract is a JSON file containing contract metadata used to interact with your contract and stored under the `lab4/build/contracts` directory. This directory is created after you compile the contracts (`npm run compile`) and modified at every new migration/deployment (`npm run migrate:ganache`).
Your client will use the [ABI](https://docs.soliditylang.org/en/v0.8.9/abi-spec.html) exported in the contract metadata to interact with the deployed contract.
