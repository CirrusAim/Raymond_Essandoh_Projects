var BettingMock = artifacts.require("BettingMock");

module.exports = function (deployer) {
    deployer.deploy(BettingMock, [web3.utils.soliditySha3("team1"), web3.utils.soliditySha3("team2")]);
};