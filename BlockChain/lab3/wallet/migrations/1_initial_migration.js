const Web3 = require('web3');
const TruffleConfig = require('../truffle-config');
const Migrations = artifacts.require('Migrations');

module.exports = function (deployer, network, accounts) {
  const config = TruffleConfig.networks[network];

  if (network == "localgeth") { // deploy on local geth
    let password = config.password();
    if (password) { // password: testCh@in
      const web3 = new Web3(new Web3.providers.HttpProvider('http://' + config.host + ':' + config.port));
      console.log("Unlocking account: " + accounts[0]);
      web3.eth.personal.unlockAccount(accounts[0], password, 600); //duration in seconds, 0 to unlock until the geth stops running.
    }
  }
  deployer.deploy(Migrations);
};
