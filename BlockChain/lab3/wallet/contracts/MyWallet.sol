// SPDX-License-Identifier: GPL-3.0-only
pragma solidity >=0.6.0 <=0.8.9;

contract MyWallet {
    address public owner; // contract owner
    uint256 internal balance; // contract balance in wei

    /* TODO (students): Take a look at the tests to figure out which
     * fields your events need to have.
     */
    
    event Deposited(address _from, uint256 _weiAmount);
    event Withdrawn(address _to, uint256 _weiAmount);
    event Transferred(address _from, address _to,uint256 _weiAmount);

    modifier onlyOwner() {
        /* TODO (students) */
        /* Should check if the sender is the contract owner */
        require(msg.sender == owner);
        _;
    }

    constructor() {
        /* TODO (students) */
        /* Should set the contract owner */
        owner = msg.sender;
    }

    /**
     * @notice Deposits the sent amount in the contract as credit to be withdrawn.
     * @dev msg.value contains the amount in Wei to be deposited.
     */
    function deposit() public payable {
        /* TODO (students) */
        // Should accept deposits from anyone
        // Should emit Deposited event
         balance += msg.value;
        // address(this).balance += msg.value;
        emit Deposited(msg.sender,msg.value);
    }

    /**
     * @notice Withdraws to the owner address the requested amount from the contract.
     * @param weiAmount The amount in Wei to be withdrawn.
     */
    function withdraw(uint256 weiAmount) public {
        /* TODO (students) */
        /* TODO (students) */
        // Should only be called by the contract owner
        if(msg.sender != owner) { revert('sender is not the owner'); }
        // Should check for sufficient funds to be withdrawn
        if(weiAmount > this.getBalance() && this.getBalance() <= 0) { revert('insufficient funds'); }
        // Should emit Withdrawn event

        balance -=  weiAmount;
        payable(msg.sender).transfer(weiAmount);
        emit Withdrawn(owner, weiAmount);
    }

    /**
     * @notice Transfers the requested amount to another address.
     * @param recipient The address of the recipient.
     * @param weiAmount The amount in Wei to be transferred.
     */
    function transfer(address payable recipient, uint256 weiAmount) public {
        /* TODO (students) */
        // Should only be called by the owner
        if(msg.sender != owner) { revert('sender is not the owner'); }
        // Should check for sufficient funds to be transferred
        if(weiAmount > this.getBalance() && this.getBalance() <= 0) { revert('insufficient funds'); }
        // Should emit Transferred event
        balance -= weiAmount;
        recipient.transfer(weiAmount);
        emit Transferred(msg.sender,recipient, weiAmount);
    }

    /**
     * @notice Returns the current balance.
     * @return The current contract balance.
     */
    function getBalance() public view returns (uint256) {
        /* TODO (students) */
        return balance;
    }
}
