const MyWalletMock = artifacts.require('MyWalletMock');
module.exports = function (deployer) {
  deployer.deploy(MyWalletMock);
};