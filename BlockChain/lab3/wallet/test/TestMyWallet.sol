// SPDX-License-Identifier: GPL-3.0-only
pragma solidity >=0.6.0 <=0.8.9;

import "truffle/Assert.sol";
import "truffle/DeployedAddresses.sol";
import "../contracts/mocks/MyWalletMock.sol";

contract TestMyWalletContract {
  MyWalletMock wallet;
  // `initialBalance` is a special Truffle variable that allows
  // funding the TestMyWallet contract after deployment.
  uint public initialBalance = 2 ether;

  // https://solidity.readthedocs.io/en/v0.6.12/contracts.html?highlight=fallback#receive-ether-function
  receive() external payable {}

  function beforeEach() public {
    wallet = new MyWalletMock();
  }

  function testInitialBalance() public {
    Assert.equal(address(this).balance, initialBalance, "The contract TestMyWallet should have a starting balance of 10 ether");
    Assert.equal(wallet.getBalance(), 0, "A new Wallet contract has a zero balance");
  }

  function testSettingAnOwnerDuringCreation() public {
    // The TestMyWallet contract is the deployer address
    Assert.equal(wallet.owner(), address(this), "The owner shouldn't be different than the deployer");
  }

  function testSettingAnOwnerUsingDeployedContract() public {
    // The msg.sender(first account address in ganache) is the contract deployer.
    wallet = MyWalletMock(DeployedAddresses.MyWalletMock());
    Assert.equal(wallet.owner(), msg.sender, "The owner shouldn't be different than the deployer");
  }

  function testValidDeposit() public {
    // Call 'deposit' with 1000 wei (it's being sent by 'TestMyWallet').
    // Wei is the default denomination in Solidity if you don't 
    // specify another unit.
    uint initBalance = address(this).balance;

    wallet.deposit{value: 1000}();

    Assert.equal(address(this).balance, initBalance - 1000 wei, "Current balance of the sender doesn't correspond to the amount deposited");
    Assert.equal(wallet.getBalance(), 1000 wei, "The balance is different than the deposited amount");
  }

  function testAccumulatedDeposits() public {
    wallet.deposit{value: 3 gwei}();
    wallet.deposit{value: 15 gwei}();
    Assert.equal(wallet.getBalance(), 18 gwei, "Balance is different than sum of the deposits");
  }

  function testWithdrawalAttemptWithNoBalance() public {
    Assert.equal(wallet.getBalance(), 0, "Balance should be 0 initially");

    (bool r, ) = address(wallet).call(abi.encodeWithSelector(wallet.withdraw.selector, 1 ether));

    Assert.isFalse(r, "Should revert due to a withdrawal attempt without funds");
  }

  function testWithdrawalByAnOwner() public {
    uint initBalance = address(this).balance;
  
    wallet.setBalance{value: 100 gwei}();
    Assert.equal(address(this).balance, initBalance - 100 gwei, "Balance of the sender before the withdrawal isn't correct");

    Assert.equal(wallet.getBalance(), 100 gwei, "Contract balance is different than the value deposited");

    (bool r, ) = address(wallet).call(abi.encodeWithSelector(wallet.withdraw.selector, 100 gwei));

    Assert.isTrue(r, "Should successfully withdrawal ether from the wallet contract");
    Assert.equal(address(this).balance, initBalance, "Balance of the sender after withdrawal should be equal to the initial balance");
  }

  function testTransferByAnOwner() public {
    uint initBalance = address(this).balance;

    address payable _to = payable(0xdCad3a6d3569DF655070DEd06cb7A1b2Ccd1D3AF);
    uint toInitBalance = _to.balance;

    wallet.setBalance{value: 100 gwei}();
    Assert.equal(address(this).balance, initBalance - 100 gwei, "Balance of the sender before the transfer isn't correct");

    Assert.equal(wallet.getBalance(), 100 gwei, "Contract balance is different than the value deposited");

    (bool r, ) = address(wallet).call(abi.encodeWithSelector(wallet.transfer.selector, _to, 100 gwei));

    Assert.isTrue(r, "Should successfully transfer ether from the wallet contract to the given account");

    Assert.equal(wallet.getBalance(), 0, "Balance at contract after the transfer should be zero");

    Assert.equal(_to.balance, toInitBalance + 100 gwei, "The transfered value isn't the expected");
  }
}