const { etherToWei, currentBalance } = require('./utils');
const { BN, expectEvent, expectRevert, constants } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');

const Betting = artifacts.require('BettingMock');

contract('Betting', accounts => {
    const [owner, oracle, winner1, winner2, loser1, loser2] = accounts;
    let betting = null;
    const team1 = web3.utils.soliditySha3("team1");
    const team2 = web3.utils.soliditySha3("team2");
    const team3 = web3.utils.soliditySha3("team3"); // invalid outcome
    const outcomes = [team1, team2]; // valid outcomes

    describe('constructor', () => {
        it('should successfully create the contract with nonempty outcomes', async () => {
            betting = await Betting.new(outcomes);
            expect(await betting.owner()).to.be.equal(owner);
            expect(await betting.oracle()).to.be.equal(constants.ZERO_ADDRESS);
            for (let i = 0; i < outcomes.length; i++) {
                expect(await betting.validOutcomes(outcomes[i])).to.be.equal(true);
            }
        });

        it('should successfully get a deployed contract', async () => {
            betting = await Betting.deployed(outcomes);
            expect(await betting.owner()).to.be.equal(owner);
            expect(await betting.oracle()).to.be.equal(constants.ZERO_ADDRESS);
            for (let i = 0; i < outcomes.length; i++) {
                expect(await betting.validOutcomes(outcomes[i])).to.be.equal(true);
            }
        });

        it('should revert when creating contract with empty outcomes', async () => {
            expectRevert(Betting.new([]), 'must register at least 2 outcomes');
        });

        it('should revert when creating contract with only one outcome', async () => {
            await expectRevert(Betting.new([team1]), 'must register at least 2 outcomes');
        });
    });

    describe('checking outcome', () => {
        before(async () => {
            betting = await Betting.new(outcomes);
        });

        it('should be able to check first registered outcomes after creation', async () => {
            expect(await betting.checkOutcome(team1)).to.be.bignumber.equal(new BN(0));
        });

        it('should be able to check second registered outcomes after creation', async () => {
            expect(await betting.checkOutcome(team2)).to.be.bignumber.equal(new BN(0));
        });

        it('should reverts when attempt to check an unregistered outcome', async () => {
            await expectRevert(betting.checkOutcome(team3), 'outcome not registered');
        });

        it('should shows the sum of the already betted outcomes', async () => {
            await betting.setBet(team1, { from: winner1, value: etherToWei('1') });

            expect(await betting.checkOutcome(team1)).to.be.bignumber.equal(new BN(etherToWei('1')));

            await betting.setBet(team1, { from: winner2, value: etherToWei('2') });
            expect(await betting.checkOutcome(team1)).to.be.bignumber.equal(new BN(etherToWei('3')));
            await betting.setBet(team2, { from: loser1, value: etherToWei('2') });
            expect(await betting.checkOutcome(team2)).to.be.bignumber.equal(new BN(etherToWei('2')));
        });
    });

    describe('choosing an oracle', () => {
        beforeEach(async () => {
            betting = await Betting.new(outcomes);
        });

        it('should allow the owner to define an oracle', async () => {
            const { logs } = await betting.chooseOracle(oracle);

            expectEvent.inLogs(logs, 'OracleChanged', {
                previousOracle: constants.ZERO_ADDRESS,
                newOracle: oracle
            });
        });

        it("should correctly checks if an oracle is defined", async () => {
            await betting.setOracle(oracle);
            expect(await betting.isOracle(oracle)).to.equal(true);
            expect(await betting.isOracle(owner)).to.equal(false);
        })

        it('should only allowed the owner to choose an oracle', async () => {
            await expectRevert(betting.chooseOracle(oracle, { from: winner1 }), "sender isn't the owner")
        });

        it('should not allow the owner to be an oracle', async () => {
            await expectRevert(betting.chooseOracle(owner), "the owner cannot be an oracle");
        });

        it('should not allow a gambler to be an oracle', async () => {
            await betting.addGamblers([winner1]);
            await expectRevert(betting.chooseOracle(winner1), "the oracle cannot be a gambler");
        });
    });

    describe('making a bet', () => {
        beforeEach(async () => {
            betting = await Betting.new(outcomes);
            await betting.setOracle(oracle);
        });

        it('should successfully create a bet and add a gambler', async () => {
            let gamblerBalance = await currentBalance(winner1);
            let amount = etherToWei('1');

            // Estimate the gas cost
            let gasPrice = new BN(await web3.eth.getGasPrice())
            let estimateGas = new BN(await betting.contract.methods.makeBet(team1).estimateGas({ from: winner1, value: amount }));
            let gasUsed = new BN(estimateGas.mul(gasPrice))

            await betting.makeBet(team1, { from: winner1, value: amount });

            let newGamblerBalance = await currentBalance(winner1);

            expect(await betting.getGamblers()).contains(winner1);
            let bet = await betting.bets(winner1);
            expect(bet.amount).to.be.bignumber.equal(amount);
            expect(bet.outcome).to.equal(team1);
            expect(newGamblerBalance).to.be.bignumber.equal(gamblerBalance.sub(amount.add(gasUsed)));
        });

        it('should emit an event when making bet', async () => {
            let amount = etherToWei('1');
            const { logs } = await betting.makeBet(team1, { from: winner1, value: amount });

            expectEvent.inLogs(logs, 'BetMade', {
                gambler: winner1,
                outcome: team1,
                amount: amount
            });
        });

        it('should not allow the owner to make a bet', async () => {
            await expectRevert(betting.makeBet(team1, { from: owner, value: etherToWei('1') }), "the owner cannot bet");
        });

        it('should not allow the oracle to make a bet', async () => {
            await expectRevert(betting.makeBet(team1, { from: oracle, value: etherToWei('1') }), "the oracle of the betting cannot bet");
        });

        it('should not allow anyone to make new bet after decision', async () => {
            await betting.setWinner(winner1, { value: etherToWei('1') });
            await expectRevert(betting.makeBet(team1, { from: loser2, value: etherToWei('1') }), "cannot bet after decision was made");
        });

        it('should not allow a gambler to bet twice', async () => {
            await betting.makeBet(team1, { from: winner1, value: etherToWei('0.4') });

            let bet = await betting.bets(winner1);
            expect(bet.outcome).to.equal(team1);
            expect(await betting.getGamblers()).contains(winner1);

            await expectRevert(betting.makeBet(team2, { from: winner1, value: etherToWei('0.7') }), "each gambler can only bet once");
        });

        it('should only allow bets on valid outcomes', async () => {
            await expectRevert(betting.makeBet(team3, { from: winner1, value: etherToWei('1') }), "outcome not registered");
        });

        it('should revert if there is no oracle assigned', async () => {
            await betting.resetOracle();
            await expectRevert(betting.makeBet(team1, { from: winner1, value: etherToWei('1') }), "no oracle found");
        });
    });

    describe('making decision', () => {
        beforeEach(async () => {
            betting = await Betting.new(outcomes);
            await betting.setOracle(oracle);
            await betting.addGamblers([winner1, loser1]);
            await betting.setBet(team1, { from: winner1, value: etherToWei('1') });
            await betting.setBet(team2, { from: loser1, value: etherToWei('2') });
        });

        it('should reverts if try to make decision for unregistered outcome', async () => {
            await expectRevert(betting.makeDecision(team3, { from: oracle }), "outcome not registered");
        });

        it('should reverts if makeDecision called by non oracle', async () => {
            await expectRevert(betting.makeDecision(team1, { from: owner }), "sender isn't the oracle");
        });

        it('should successfully make a decision', async () => {
            await betting.makeDecision(team1, { from: oracle });

            expect(await betting.wins(winner1)).to.be.bignumber.equal(etherToWei('3'));

            expect(await betting.wins(loser1)).to.be.bignumber.equal(etherToWei('0'));
        });

        // TODO: check make decision event emission.

        it('should not decide twice', async () => {
            await betting.makeDecision(team1, { from: oracle });
            await expectRevert(betting.makeDecision(team2, { from: oracle }), 'can make decision only once');
        });
    });

    describe("withdraw", () => {
        beforeEach(async () => {
            betting = await Betting.new(outcomes);
            await betting.setWinner(winner1, { value: etherToWei('1') });
        });

        it('should allow winners to withdraw their reward', async () => {
            let amount = etherToWei('1');
            let initialWinnerBalance = await currentBalance(winner1);
            let contractBalance = await currentBalance(betting.address);
            expect(contractBalance).to.be.bignumber.equal(amount);

            await betting.withdraw(amount, { from: winner1 });

            expect(await currentBalance(betting.address)).to.be.bignumber.equal(new BN(0));
            expect(await currentBalance(winner1)).to.be.bignumber.greaterThan(initialWinnerBalance);
        });

        it('should emits an event when successfully withdraw a reward', async () => {
            let amount = etherToWei('1');
            const { logs } = await betting.withdraw(amount, { from: winner1 });

            expectEvent.inLogs(logs, 'Withdrawn', {
                gambler: winner1,
                amount: amount
            });
        });

        it('should reverts when try to withdraw an amount bigger than the reward', async () => {
            await expectRevert(betting.withdraw(etherToWei('3'), { from: winner1 }), 'insufficient requested amount');
        });

        it('should only allow winners to withdraw', async () => {
            await expectRevert(betting.withdraw(etherToWei('1'), { from: loser1 }), 'sender should be a winner');
        });
    });

    describe("check winnings", () => {
        beforeEach(async () => {
            betting = await Betting.new(outcomes);
            await betting.setWinner(winner1, { value: etherToWei('3') });
        });

        it('should allow winners to check their reward', async () => {
            expect(await betting.checkWinnings({ from: winner1 })).to.be.bignumber.equal(etherToWei('3'));
        });

        it('should return 0 for anyone that isn\'t a winner', async () => {
            expect(await betting.checkWinnings({ from: loser1 })).to.be.bignumber.equal(etherToWei('0'));
        });
    });

    describe('proportional rewards for gamblers', () => {
        beforeEach(async () => {
            betting = await Betting.new([team1, team2, team3]);
            await betting.setOracle(oracle);
            await betting.addGamblers([winner1, winner2, loser1, loser2]);
            await betting.setBet(team1, { from: winner1, value: etherToWei('1') });
            await betting.setBet(team1, { from: winner2, value: etherToWei('2') });
            await betting.setBet(team2, { from: loser1, value: etherToWei('2') });
            await betting.setBet(team3, { from: loser2, value: etherToWei('1') });
        });

        it('winners should get proportional reward', async () => {
            await betting.makeDecision(team1, { from: oracle });
            expect(await betting.wins(winner1)).to.be.bignumber.equal(etherToWei('2'));
            expect(await betting.wins(winner2)).to.be.bignumber.equal(etherToWei('4'));
            expect(await betting.wins(loser1)).to.be.bignumber.equal(etherToWei('0'));
            expect(await betting.wins(loser2)).to.be.bignumber.equal(etherToWei('0'));
        });
    });

    describe('no gamblers bet on the correct outcome, the oracle wins', () => {
        before(async () => {
            betting = await Betting.new([team1, team2, team3]);
            await betting.setOracle(oracle);
            await betting.addGamblers([loser1, loser2]);
            await betting.setBet(team1, { from: loser1, value: etherToWei('1') });
            await betting.setBet(team2, { from: loser2, value: etherToWei('2') });
        });

        it('winners get proportional reward', async () => {
            await betting.makeDecision(team3, { from: oracle });

            expect(await betting.wins(oracle)).to.be.bignumber.equal(etherToWei('3'));
            expect(await betting.wins(loser1)).to.be.bignumber.equal(etherToWei('0'));
            expect(await betting.wins(loser2)).to.be.bignumber.equal(etherToWei('0'));
        });
    });

    describe('reseting contract state', () => {
        beforeEach(async () => {
            betting = await Betting.new([team1, team2, team3]);
            await betting.chooseOracle(oracle);
            await betting.makeBet(team1, { from: winner1, value: etherToWei('1') });
            await betting.makeBet(team1, { from: winner2, value: etherToWei('2') });
            await betting.makeBet(team2, { from: loser1, value: etherToWei('2') });
            await betting.makeBet(team3, { from: loser2, value: etherToWei('1') });
        });

        it('should not allow reset before decision', async () => {
            await expectRevert(betting.contractReset(), "cannot reset before decision");
        });

        it('should allow reset after decision', async () => {
            await betting.makeDecision(team1, { from: oracle });
            await betting.contractReset();
            expect(await betting.owner(), owner, "owner should not be reseted");
            expect(await betting.oracle(), oracle, "oracle doesn't need to be reset");
            expect(await betting.decisionMade()).to.equal(false);
            expect((await betting.getGamblers()).length).to.equal(0);
            expect((await betting.getOutcomes()).length).to.equal(0);
        });
    });
});
