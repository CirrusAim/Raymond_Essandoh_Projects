const { toWei, fromWei } = require('web3-utils');
const { BN } = require('@openzeppelin/test-helpers');

etherToWei = (amount) => {
    return new BN(toWei(amount, 'ether'));
};

weiToEther = (amount) => {
    return fromWei(amount, 'ether');
};

currentBalance = async (address) => {
    return new BN(await web3.eth.getBalance(address));
};

module.exports = {
    currentBalance,
    etherToWei,
    weiToEther
};