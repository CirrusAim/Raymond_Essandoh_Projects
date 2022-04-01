// SPDX-License-Identifier: GPL-3.0-only
pragma solidity >=0.6.0 <=0.8.9;

import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "../MyWallet.sol";

contract MyWalletMock is MyWallet {
    using SafeMath for uint256;

    function setBalance() public payable {
        // NOTE: `SafeMath` is generally not needed starting with Solidity 0.8,
        // since the compiler now has built in overflow checking.
        // But it still used here since we support solidity version >= 0.6.
        balance = balance.add(msg.value);
    }
}
