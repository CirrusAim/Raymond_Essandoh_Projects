# Lab 4: Single-decree Paxos and Multi-Paxos

| Lab 4: | Single-decree Paxos and Multi-Paxos |
| ---------------------    | --------------------- |
| Subject:                 | DAT520 Distributed Systems |
| Deadline:                | **March 21, 2022 23:59** |
| Expected effort:         | 40-60 hours |
| Grading:                 | Graded |
| Submission:              | Group |

## Table of Contents

1. [Introduction](#introduction)
2. [Resources](#resources)
3. [Questions](#questions)
4. [Paxos](#paxos)
5. [Multi-Paxos](#multi-paxos)

## Introduction

The overall objective of this lab is to implement a single-decree and a multi-decree version of Paxos (also known as Multi-Paxos).
The assignment consist of three parts:

1. A set of theory questions that you should answer.

2. Implementation of the single-decree Paxos algorithm for each of the three Paxos roles: Proposer, Acceptor and Learner, as described [here](singlepaxos/README.md).
   This implementation has corresponding unit tests and will be automatically verified by QuickFeed.

3. Implementation of the multi-decree Paxos algorithm for each of the three Paxos roles as described [here](multipaxos/README.md).
   This implementation has corresponding unit tests and will be automatically verified by QuickFeed.

## Resources

Several Paxos resources are listed below. You should use these resources to
answer the [questions](#questions) for this lab. You are also advised to use
them as support literature when working on your implementation now and in
future lab assignments.

* [Paxos Explained from Scratch](resources/paxos-scratch-slides.pdf)  - slides.
* [Paxos Explained from Scratch ](resources/paxos-scratch-paper.pdf) - paper.
* [Paxos Made Insanely Simple](resources/paxos-insanely-simple.pdf) - slides. Also
  contains pseudo code for Proposer and Acceptor.
* [Paxos Made Simple](resources/paxos-simple.pdf)
* [Paxos Made Moderately Complex](resources/paxos-made-moderately-complex.pdf)
* [Paxos Made Moderately Complex (ACM Computing Surveys)](resources/a42-renesse.pdf)
* [Paxos for System Builders](resources/paxos-system-builders.pdf)
* [The Part-time Parliament](resources/part-time-parliment.pdf)

## Questions

Answer the questions below. You should write down and submit your answers by
using the [answers.md](answers.md) file.

1. Is it possible that Paxos enters an infinite loop? Explain.

2. Is the value to agree on included in the `Prepare` message?

3. Does Paxos rely on an increasing proposal/round number in order to work?
   Explain.

4. Look at this description for Phase 1B: If the proposal number N is larger
   than any previous proposal, then each Acceptor promises not to accept
   proposals less than N, and sends the value it last accepted for this
   instance to the Proposer. What is meant by "the value it last accepted"? And
   what is an "instance" in this case?

5. Explain, with an example, what will happen if there are multiple
   proposers.
   
6. What happens if two proposers both believe themselves to be the leader and
   send `Prepare` messages simultaneously?

7. What can we say about system synchrony if there are multiple proposers (or
   leaders)?

8. Can an acceptor accept one value in round 1 and another value in round 2?
   Explain.

9. What must happen for a value to be "chosen"? What is the connection between
   chosen values and learned values?

## Paxos

The first Paxos implementation for this lab will be a single-decree version as described [here](singlepaxos/README.md).
This variant of Paxos is only expected to choose a single command. It does not need
several slots.

## Multi-Paxos

After implementing the core logic of the single-decree Paxos, your next task is to implement its variation called Multi-Paxos as described [here](multipaxos/README.md), which will allow you to choose multiple commands.
