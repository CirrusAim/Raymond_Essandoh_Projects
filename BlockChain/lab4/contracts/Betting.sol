// SPDX-License-Identifier: GPL-3.0-only
pragma solidity >=0.6.0 <=0.8.10;

contract Betting {
    /* Define the Bet struct */
    struct Bet {
        bytes32 outcome; // the guessed outcome
        uint256 amount; // the bet amount
    }

    address public owner; // the contract owner
    address public oracle; // the oracle that will decide the outcome of the betting
    address[] public gamblers; // list of the gamblers addresses
    mapping(address => bool) public isGambler;

    /* Maps the gambler's addresses to their bets */
    mapping(address => Bet) public bets;

    bool public decisionMade;

    /* List of all winners (maps are not iterable in solidity)*/
    address[] public winners;
    /* Maps the winners to their prize amount */
    mapping(address => uint256) public wins;

    /* Keep track of all possible outcomes */
    bytes32[] public outcomes;
    /* Map the valid outcomes for his betting */
    mapping(bytes32 => bool) public validOutcomes;
    /* Keep track of the total bet amount in all outcomes */
    /* How much amount a user has bet on all teams(outcomes) */
    mapping(bytes32 => uint256) public outcomeBets;

    /* Add any events you think are necessary */
    event BetMade(
        address indexed gambler,
        bytes32 indexed outcome,
        uint256 amount
    );
    event Winners(address[] indexed wins, uint256 totalPrize);
    event OracleChanged(
        address indexed previousOracle,
        address indexed newOracle
    );
    event Withdrawn(address indexed gambler, uint256 amount);

    /* Uh Oh, what are these? */
    modifier onlyOwner() {
        require(msg.sender == owner, "sender isn't the owner");
        _;
    }

    modifier onlyOracle() {
        require(isOracle(msg.sender), "sender isn't the oracle");
        _;
    }

    modifier requireOracle() {
        require(oracle != address(0), "no oracle found");
        _;
    }

    modifier outcomeExists(bytes32 outcome) {
        require(validOutcomes[outcome], "outcome not registered");
        _;
    }

    modifier onlyWinners() {
        require(wins[msg.sender] > 0, "sender should be a winner");
        _;
    }

    /* Constructor function, where owner and outcomes are set */
    constructor(bytes32[] memory initOutcomes) {
        // should register at least 2 possible outcomes,
        if (initOutcomes.length < 2){
            revert("must register at least 2 outcomes");
        }else{

            setOutcomes(initOutcomes);
        }
        // and define the contract owner.
        owner = msg.sender;

    }

    function setOutcomes(bytes32[] memory _outcomes) public {
        // I guess this function can be used to initialize the outcomes
            uint i = 0;
            for (i = 0; i < _outcomes.length;i++){
                outcomes.push(_outcomes[i]);
                validOutcomes[_outcomes[i]] = true;
            }
    }

    /**
     * @notice This function allows owner to chooses their trusted Oracle.
     * @param newOracle The address of the new oracle.
     */
    function chooseOracle(address newOracle) public onlyOwner(){
        // Must be called only by the contract owner
        // The oracle cannot be neither a gambler or the owner
        require(newOracle != owner, "the owner cannot be an oracle" );
        require(isGambler[newOracle] == false, "the oracle cannot be a gambler");

        // Should emit OracleChanged event
        address prevOrac = oracle;
        oracle = newOracle;
        emit OracleChanged(prevOrac, newOracle);
    }

    /**
     * @notice Make a bet.
     * @param outcome The hash of the outcome to bet on.
     */
    function makeBet(bytes32 outcome) public  requireOracle() outcomeExists(outcome) payable {
        // Owner and oracle cannot make a bet
        // An oracle should be assigned before starting bets
        // Must be impossible to bet after decision was made
        // Betters are registered by placing a bet
        // A gambler can only bet on a registered outcome
        require(msg.sender != owner, "the owner cannot bet");
        require(msg.sender != oracle, "the oracle of the betting cannot bet");
        require(decisionMade == false, "cannot bet after decision was made");
        
        // A gambler cannot bet twice
        require(isGambler[msg.sender] == false,"each gambler can only bet once");


        isGambler[msg.sender] = true;      
        gamblers.push(msg.sender);
        // bets[msg.sender] = Bet();
        bets[msg.sender] = Bet({outcome:outcome, amount: msg.value});
        // Should emit BetMade event
        emit BetMade(msg.sender, outcome, msg.value);
    }

    /**
     * @notice Decide on an outcome.
     * @param decidedOutcome The chosen outcome.
     */
    function makeDecision(bytes32 decidedOutcome) public onlyOracle() outcomeExists(decidedOutcome){
        /* TODO (students) */
        // Must be called only by the oracle
        
        // The oracle must chooses which outcome wins calling and set the winners.
        // Should be called only once before reset
        require(decisionMade == false, "can make decision only once");
        // Winning outcome must exist
        

        // Gamblers and bets must exists before make a decision
        require(gamblers.length > 0, "No gamblers exists");

        uint winCount = 0;
        uint256 totalPrize  = 0;
        uint256 totalAfterWinnerBetReverted = 0;
        uint256 totalBetReverted = 0;
        for(uint i = 0 ; i< gamblers.length; i++){
            totalPrize += bets[gamblers[i]].amount;
            if (bets[gamblers[i]].outcome == decidedOutcome){
                winners.push(gamblers[i]);
                wins[gamblers[i]] = bets[gamblers[i]].amount; 
                winCount += 1;
                totalBetReverted += bets[gamblers[i]].amount; 
            }else{
                totalAfterWinnerBetReverted += bets[gamblers[i]].amount;
            }
        }

        if (winCount == 0){
            // Oracle wins
            winners.push(oracle);
            wins[oracle] = totalPrize;
        }else{
            
            uint256 winnerAmount = (totalAfterWinnerBetReverted)/totalBetReverted;
            for (uint i=0;i< winners.length;i++){
                wins[winners[i]] += (winnerAmount * bets[gamblers[i]].amount);
            }
        }
        
        // The winners receive a proportional share of the total funds at stake if they all bet on the correct outcome
        // If all gamblers bet on the correct outcome, then they must get reimbursed their funds.
        // If no gamblers bet on the correct outcome, then the oracle wins the sum of the funds.
        // Should emit Winners event
        decisionMade = true;
        emit Winners(winners, totalPrize);
    }

    /**
     * @notice This function allows the winners to withdraw their
     * winnings safely (if they win something).
     * @param amount The amount to be withdrawn
     */
    function withdraw(uint256 amount) public onlyWinners(){
        /* TODO (students) */
        // Should only be called by winners
        // Winners can withdraw multiple times until the total amount of their prize
        require(amount <= wins[msg.sender],"insufficient requested amount");
        wins[msg.sender] -= amount; 
        address payable th = payable(msg.sender);
        th.transfer(amount);
        emit Withdrawn(msg.sender, amount);
    }

    /**
     * @notice Reset the contract state
     */
    function contractReset() public onlyOwner(){
        // Must only be called by the contract owner.
        // Should not allow reset the contract state before a decision is made
        require(decisionMade == true, "cannot reset before decision");
        // Must reset the contract variables to the initial state to allow new bettings and outcomes.
        delete gamblers;
        delete outcomes;
        delete winners;
        decisionMade = false;
    }

    /**
     * @dev This function allows anyone to check the amount
     * already betted per outcome.
     * @param outcome The hash of the outcome to be checked.
     * @return The amount betted for the given outcome.
     */
    function checkOutcome(bytes32 outcome) public outcomeExists(outcome) view returns (uint256) {
        // Must revert if outcome does not exist
        // Returns the current stake for the given outcome
        return outcomeBets[outcome];
    }

    /**
     * @notice This function is similar to `checkOutcome` but
     * it receives the outcome as string.
     * @param outcomeString The string representation of the outcome
     * to be checked.
     * @dev It uses the `keccak256` hash function to get the hash of the `outcomeString`
     * @return The amount betted for the given outcome.
     */
    function checkOutcomeString(string memory outcomeString)
        public
        view
        returns (uint256)
    {
        // hash outcomestring using keccak256. Then checkOutcome
        bytes32 add = keccak256(abi.encodePacked(outcomeString));
        return outcomeBets[add];
    }

    /**
     * @notice This function allows anyone to check their winnings.
     * @return The winning amount of the msg.sender if it exists.
     */
    function checkWinnings() public view returns (uint256) {
        return wins[msg.sender];
    }

    function isOracle(address _oracle) public view returns (bool) {
        return oracle == _oracle;
    }

    function getGamblers() public view returns (address[] memory) {
        return gamblers;
    }

    function getWinners() public view returns (address[] memory) {
        return winners;
    }

    function getOutcomes() public view returns (bytes32[] memory) {
        return outcomes;
    }
}
