const { etherToWei, currentBalance } = require('./utils');
const { BN, expectEvent, expectRevert } = require('@openzeppelin/test-helpers');
const { expect } = require('chai');

const MyWalletMock = artifacts.require('MyWalletMock');

contract('MyWallet', accounts => {
    const [owner, recipient, giver] = accounts;
    let wallet = null;

    describe('constructor', () => {
        it('should successfully create the contract', async () => {
            wallet = await MyWalletMock.new(); //deploy a new contract
            expect(await wallet.owner()).to.be.equal(owner);
            expect(await wallet.getBalance()).to.be.bignumber.equal(new BN(0));
        });

        it('should successfully get a deployed contract', async () => {
            wallet = await MyWalletMock.deployed(); // contract deployed by migration
            expect(await wallet.owner()).to.be.equal(owner);
            expect(await wallet.getBalance()).to.be.bignumber.equal(new BN(0));
        });
    });

    describe('deposit', () => {
        beforeEach(async () => {
            wallet = await MyWalletMock.new();
        });

        it('should successfully make a deposit', async () => {
            let amount = etherToWei('1');
            await wallet.deposit({ value: amount });
            expect(await wallet.getBalance()).to.be.bignumber.equal(amount);
        });

        it('should emit an event when create a new deposit', async () => {
            let amount = etherToWei('1');
            const { logs } = await wallet.deposit({ from: owner, value: amount });

            expectEvent.inLogs(logs, 'Deposited', {
                _from: owner,
                _weiAmount: etherToWei('1') // 1000000000000000000 wei == 1 ether
            });
        });

        it('should allow anyone to make a deposit', async () => {
            let giverBalance = await currentBalance(giver);
            let amount = etherToWei('5');

            // Estimate the gas cost
            let gasPrice = new BN(await web3.eth.getGasPrice())
            let estimateGas = new BN(await wallet.contract.methods.deposit().estimateGas({ from: giver, value: amount }));
            let gasUsed = new BN(estimateGas.mul(gasPrice))

            const { logs } = await wallet.deposit({ from: giver, value: amount });

            expectEvent.inLogs(logs, 'Deposited', {
                _from: giver,
                _weiAmount: amount
            });

            let newGiverBalance = await currentBalance(giver);
            expect(await wallet.getBalance()).to.be.bignumber.equal(amount);
            expect(newGiverBalance).to.be.bignumber.equal(giverBalance.sub(amount.add(gasUsed)));
        });
    });

    describe('withdraw', () => {
        beforeEach(async () => {
            wallet = await MyWalletMock.new();
        });

        it('should revert from a withdraw attempt with no funds', async () => {
            await expectRevert(wallet.withdraw(etherToWei('1')), 'insufficient funds');
        });

        it('should revert from an unauthorized withdraw attempt', async () => {
            await wallet.setBalance({ value: etherToWei('2') });
            expect(await wallet.getBalance()).to.be.bignumber.equal(etherToWei('2'));

            await expectRevert(wallet.withdraw(etherToWei('1'), { from: giver }), 'sender is not the owner');
        });

        it('should successfully make a withdraw', async () => {
            await wallet.setBalance({ value: etherToWei('3') });
            expect(await wallet.getBalance()).to.be.bignumber.equal(etherToWei('3'));

            let initialBalance = await currentBalance(owner);
            let amount = etherToWei('2');
            const { logs } = await wallet.withdraw(amount);

            expectEvent.inLogs(logs, 'Withdrawn', {
                _to: owner,
                _weiAmount: amount
            });

            expect(await wallet.getBalance()).to.be.bignumber.equal(etherToWei('1'));
            expect(await currentBalance(owner)).to.be.bignumber.greaterThan(initialBalance);
        });
    });

    describe('transfer', () => {
        beforeEach(async () => {
            wallet = await MyWalletMock.new();
        });

        it('should revert from a transfer attempt with no funds', async () => {
            await expectRevert(wallet.transfer(recipient, etherToWei('1')), 'insufficient funds');
        });

        it('should revert from an unauthorized transfer attempt', async () => {
            await wallet.setBalance({ value: etherToWei('2') });
            expect(await wallet.getBalance()).to.be.bignumber.equal(etherToWei('2'));

            await expectRevert(wallet.transfer(recipient, etherToWei('1'), { from: giver }), 'sender is not the owner');
        });

        it('should successfully make a transfer', async () => {
            await wallet.setBalance({ value: etherToWei('3') });
            expect(await wallet.getBalance()).to.be.bignumber.equal(etherToWei('3'));

            let initialBalance = await currentBalance(recipient);
            let amount = etherToWei('2');
            const { logs } = await wallet.transfer(recipient, amount);

            expectEvent.inLogs(logs, 'Transferred', {
                _from: owner,
                _to: recipient,
                _weiAmount: amount
            });

            expect(await wallet.getBalance()).to.be.bignumber.equal(etherToWei('1'));
            expect(await currentBalance(recipient)).to.be.bignumber.greaterThan(initialBalance);
        });
    });
});
