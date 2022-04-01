const Betting = artifacts.require('Betting');
module.exports = function (deployer) {
  deployer.deploy(Betting,[web3.utils.soliditySha3("Pakistan"), web3.utils.soliditySha3("England"),web3.utils.soliditySha3("Australia"),web3.utils.soliditySha3("New Zealand")]);
};
