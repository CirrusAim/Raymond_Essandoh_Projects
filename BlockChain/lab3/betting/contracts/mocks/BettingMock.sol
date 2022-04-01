// SPDX-License-Identifier: GPL-3.0-only
pragma solidity >=0.6.0 <=0.8.9;

import "../Betting.sol";

contract BettingMock is Betting {
    constructor(bytes32[] memory initOutcomes) Betting(initOutcomes) {}

    function addGamblers(address[] memory _gamblers) public {
        for (uint256 i = 0; i < _gamblers.length; i++) {
            gamblers.push(_gamblers[i]);
            isGambler[_gamblers[i]] = true;
        }
    }

    function setOracle(address _oracle) public {
        oracle = _oracle;
    }

    function resetOracle() public {
        oracle = address(0);
    }

    function setBet(bytes32 outcome) public payable {
        bets[msg.sender] = Bet(outcome, msg.value);
        outcomeBets[outcome] += msg.value;
    }

    function setWinner(address winner) public payable {
        decisionMade = true;
        winners.push(winner);
        wins[winner] = msg.value;
    }
}
