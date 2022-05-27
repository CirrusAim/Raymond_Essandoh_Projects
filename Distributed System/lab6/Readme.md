# Lab 6: Bank Application with Reconfiguration

| Lab 6: | Bank Application with Reconfiguration |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **April 29, 2022 23:59** |
| Expected effort:         | 30-40 hours |
| Grading:                 | Graded |
| Submission:              | Group |

## Table of Contents

1. [Introduction](#introduction)
2. [Task 1: Bank Application](#task-1-bank-application)
3. [Task 2: Dynamic Membership through Reconfiguration](#task-2-dynamic-membership-through-reconfiguration)
4. [References](#references)

## Introduction
The overall objective of this lab is to implement a resilient bank application that stores a set of bank accounts and apply transactions to them as they are decided by your gorums Multi-paxos nodes from the [lab5](../lab5/gorumspaxos/README.md). The assignment consists of two parts:

1. You will use your implementation from previous lab to replicate a set of bank accounts. You will also be required to extend the [client](../lab5/gorumspaxos/cmd/paxosclient/main.go) to handle messages to your bank application running on the replicas.

2. You will extend your implementation to enable dynamic membership into your Multi-Paxos protocol, allowing your system to reconfigure the set of nodes available and to keep running in the presence of failures of some nodes.

**Note** that no part of this lab will be verified by QuickFeed, and they will be verified by a member of the teaching staff during lab hours.

## Task 1: Bank Application

You will need to extend your replica from lab5 and add the following functionalities:

* Store bank accounts.
* Apply bank transactions in correct order as they are decided by your Multi-Paxos nodes.

Your system should in this assignment use the Multi-Paxos protocol to replicate a set of bank accounts and their balances information.
Clients can issue transactions to the accounts.
Transactions sent by clients should be processed in the same order by all replicas using the Multi-Paxos protocol you have implemented.

Thus, you will need to define an account type and keep a `map` or a list of accounts balances per account number or similar data structure as shown in the example below. You also need to define and implement the bank operations over the account type depending on the type of transactions received.
The operation should be executed in a `Process` method of an account. This method performs the correspondent action depending of the `tx.Op`.

There are three basic operations that you must implement. But only two of them modify the accounts state:
- `Deposit`: adds the `tx.Amount` to the account balance or return the error message `Can't deposit negative amount (X NOK)` when a negative amount `X` is given.
- `Withdrawal`: subtracts the `tx.Amount` from the account balance or return the error message `Not enough funds for withdrawal. Balance: X NOK - Requested Y NOK` if `X < Y`.
- `Balance`: this operation does not modify the state and should just return the current account balance.

* Note: you can add more operation and error messages if you wish. How you will design the bank application is up to you. You can also create a Bank struct and use it in the replica or just implement the bank logic in the Replica itself like illustrated below. In either way, your bank application must fullfil the functionalities described in this document.

```go
type PaxosReplica struct {
	// ...
	accounts map[uint32]*Account
}

// Account represents a bank account with an account number and a balance.
type Account struct {
	Number  uint32
	Balance int32
}

// Process applies transaction tx to Account a and returns a corresponding TransactionResult
func (a *Account) Process(tx *pb.Transaction) *pb.TransactionResult
```

A `Transaction` is sent by clients together with an account number in the `Value` message.
The actual value is not anymore a `clientCommand` as in the previous lab, but it is now a combination of the `accountNumber` field and the `tx` field of type `Transaction`.

A transaction has a operation field which defines which kind of operation the bank application running on the replica should perform, and an optional amount field.

Thus you should also change your implementation for Multi-Paxos in a way to be able to handle bank transactions.
The `Value` type definition in the lab5 [proto files](../lab5/gorumspaxos/proto/multipaxos.proto) must be changed to the following:

```proto
// Transaction represents a bank transaction with an operation type and amount.
// If Op == Balance then the Amount field should be ignored.
// Basic operations types are:
// - Balance == 0
// - Deposit == 1
// - Withdrawal == 2
// You may add more if you like.
message Transaction {
	int32 Op = 1;
	int32 amount = 2;
}

message Value {
	string clientID = 1;
	uint32 clientSeq = 2;
	bool isNoop = 3;
	uint32 accountNumber = 4;
	Transaction tx = 5;
}
```

The method should be used by a replica when applying a decided transaction to an account.
The resulting `TransactionResult` should be used by a replica in the reply to the client.

In Lab 5 the reply to a client was simply the value it originally sent (i.e. `resp.ClientCommand == req.ClientCommand`).
However, for this lab the reply is the result of execute a client's transaction over the current replica's state (i.e. bank account balances).
The `Response` should therefore be modified in the proto file as well, to use the transaction result instead of a `clientCommand`.
The `ClientID`, `ClientSeq` and `AccountNumber` fields from the response should be populated from the corresponding decided value.

```proto
// TransactionResult represents a result of processing a Transaction for an
// account with AccountNumber. Any error processing the transaction is reported
// in the ErrorMessage field. If an error is reported then the Balance field
// should be ignored.
message TransactionResult {
	uint32 accountNumber = 1;
	int32 balance = 2;
	string errorMessage = 3;
}

message Response {
	string clientID = 1;
	uint32 clientSeq = 2;
	TransactionResult txResult = 3;
}
```

A replica should generate a response after execute a transaction, applying the operation `tx.Op` to an account when needed.
A response should be forwarded to the client handling part of your application and sent back to the appropriate client.

You will need to modify at least the following files/methods:
- [defs.go](../lab5/gorumspaxos/defs.go) adjusting the `ClientCommand` to the new value field.
- [qspec.go](../lab5/gorumspaxos/qspecs.go) adjust the `ClientHandleQF` to handle the proper replies and validate errors.
- [replica.go](../lab5/gorumspaxos/replica.go) adjust it to store accounts information and process commited replies.
- [paxosclient](../lab5/gorumspaxos/cmd/paxosclient/main.go) to send bank transactions and print responses.
- [multipaxos.proto](../lab5/gorumspaxos/proto/multipaxos.proto) modifying the `Value` and `Response` messages and adding the `Transaction` and `TransactionResult` messages.

You may need to modify other files as well depending of your lab 5 implementation.

Furthermore, a learner may deliver decided values out of order.
The replica needs to ensure that decided values (i.e. transactions) are processed in order.
So it needs to keep track of the ID for the highest decided slot to accomplish this.
This assignment, will use the name `adu` (_all decided up to_) when referring to this variable.
It should initially be set to `-1`.

The replica should buffer out-of-order decided values and only apply them (i.e. process) when it has consecutive sequence of decided slots.
More specifically the replica should handle decided values from the learner equivalently to the logic in the following pseudo code:

```go
on receive decided value v from learner:
	handleDecideValue(v)
```

```go
handleDecidedValue(value):
	if slot id for value is larger than adu+1:
		buffer value
		return
	if value is not a noop:
		if account for account number in value is not found:
			create and store new account with a balance of zero
		apply transaction from value to account if possible (e.g. user has balance)
		create response with appropriate transaction result, client id and client seq
		forward response to client handling module
	increment adu by 1
	increment decided slot for proposer
	if has previously buffered value for (adu+1):
		handleDecidedValue(value from slot adu+1)
```

**Note**: The replica should not apply a transaction if the decided value has its `Noop` field set to true. For a noop it should only increment its `adu` and call the `IncrementAllDecidedUpTo()` method on the Proposer.

You should modify your client to allow send a bank transaction to the system by asking the user for the necessary input: account number, transaction type and amount.

## Task 2: Dynamic Membership through Reconfiguration

This task consists of adding dynamic membership into your Multi-Paxos protocol, namely, implement a reconfiguration command to adjust the set of servers that are executing in the system, e.g., adding or replacing nodes while the system is running.
Therefore, a configuration consists of a set of servers that are executing the Paxos protocol.

The Paxos protocol assumes a fixed configuration. Thus we must ensure that a _single configuration_ executes each instance of consensus.
Reconfiguration can be achieved in different ways.
You could use different configurations for different instances or enabling mechanisms to stop the current execution and perform a migration to the new configuration.

For example, you could implement reconfiguration using different configurations per instance by adding a new special command that defines the set of nodes configured for future instances.
Nodes could agree on this set as part of the consensus, in the same way that they agree for other messages, and once an agreement was reached about the new proposed configuration, all nodes could migrate to such configuration.
Your implementation should also ensure that previous instances that did not get decided were garbage collected, for example, using a `noop` command to force the instances to finalize.

Another alternative would be to prepare your system to migrate its state in case of a reconfiguration, such as creating a snapshot of the system.
In case of adding a new server at slot `i`, you could stop the consensus for future slots, ensure that you finish executing every lower-numbered slot, obtain the state immediately following execution of slot `i − 1`, then transfer this latest state to the new server (including application state), and then start the consensus algorithm again in the new configuration.

Please read section 2 of this paper [[1]](#references) for more details of those mechanisms and their advantages and problems.

You should be able to trigger a reconfiguration by sending a special
`reconfiguration request` from the bank client. You will be asked to demonstrate
different types of reconfiguration during the lab approval, e.g., simulating node failures.
Your system should be able to scale from 3 Paxos servers (which handles one failure) to 5 (2 fails) and 7 (3 fails).

* Note that this task is not as simple as it may first seem. We recommend you to study the relevant parts described here [[1]](#references) for guidance and also to take a look at an overview of the technique here [[2]](#references).

## References

1. Leslie Lamport, Dahlia Malkhi, and Lidong Zhou. _Reconfiguring a state
   machine._ SIGACT News, 41(1):63–73, March 2010.
   [[pdf]](resources/reconfig10.pdf)

2. Jacob R. Lorch, Atul Adya, William J. Bolosky, Ronnie Chaiken, John R.
   Douceur, and Jon Howell. _The smart way to migrate replicated stateful
   services._ In Proceedings of the 1st ACM SIGOPS/EuroSys European Conference
   on Computer Systems 2006, EuroSys ’06, pages 103–115, New York, NY, USA, 2006. ACM.
   [[pdf]](resources/smart06.pdf)