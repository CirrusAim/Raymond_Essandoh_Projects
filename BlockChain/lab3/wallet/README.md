# Wallet contract

Build a simple wallet contract in Solidity that stores the owner's funds.
A stub [contract](contracts/MyWallet.sol) is provided for you.
Your implementation needs to fit the following rules:

## Rules
* There can be only one contract owner.
* Anyone should be able to make a deposit.
* It should be possible for anyone to retrieve the current balance stored in the contract.
* The contract should emit events on all performed operations except `getBalance`.
* It should be possible for the owner to withdraw any amount from the current contract's balance. And **only** the owner should be authorized to do so.
* The owner should be able to transfer some amount of the contract's balance to a given address.
* If the owner attempts to withdraw or transfer more than the current contract's balance the operation should revert with the error message: `insufficient funds`.
* If someone else than the owner attempts to withdraw or transfer funds from the contract, the operation should revert with error message: `sender is not the owner`.
* You may add as many auxiliary functions or variables as you need.