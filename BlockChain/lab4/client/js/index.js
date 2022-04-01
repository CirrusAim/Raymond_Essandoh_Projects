import "./scss/index.scss";
import Web3 from "web3";
import { fromWei, toWei } from "web3-utils";
import MyWalletContract from "../../build/contracts/MyWallet.json";

const Status = Object.freeze({
    Success: "success",
    Fail: "danger",
    Info: "secondary",
    Warning: "warning"
});

const App = {
    provider: null,
    accounts: null,
    contractOwner: null,
    selectedAccount: null,
    contract: null,

    start: async function () {
        const { provider } = this;

        try {
            // get contract instance
            const networkId = await provider.eth.net.getId();
            const deployedNetwork = MyWalletContract.networks[networkId];
            this.contract = new provider.eth.Contract(
                MyWalletContract.abi,
                deployedNetwork.address,
            );

            // get accounts
            this.accounts = await this.getAccounts();
            this.setDefaultAccount();
            this.buildAccountList();
            await this.refreshPageInfo();
        } catch (error) {
            console.error(`Could not connect to contract or chain. Due the following error: ${error}`);
        }
    },

    deposit: async function () {
        const amount = document.getElementById("depositAmount").value;

        try {
            if (!this.isValidInputValue(amount)) {
                throw new Error("Amount should be greater than 0");
            }
            const amountWei = toWei(amount, "wei");
            const { deposit } = this.contract.methods;
            deposit().send({ from: this.selectedAccount, value: amountWei })
                .once("transactionHash", (txHash) => {
                    this.setStatus(Status.Info, `Transaction ${txHash} sent.`);
                })
                .once("receipt", (receipt) => {
                    // You can also monitor or inspect the event logs:
                    // https://web3js.readthedocs.io/en/v1.2.0/web3-eth-contract.html#id38
                    this.setStatus(Status.Success, `Transaction complete! Successfully deposit ${amountWei} wei into the wallet contract ${receipt.to}.`);
                })
                .on("error", (err) => {
                    this.setStatus(Status.Fail, `Transaction failed: ${err}`);
                });
        } catch (err) {
            this.setStatus(Status.Fail, `Transaction failed: ${err.message}`);
        }
        await this.refreshPageInfo();
    },

    withdraw: async function () {
        const amount = document.getElementById("withdrawAmount").value;

        try {
            if (!this.isValidInputValue(amount)) {
                throw new Error("Amount should be greater than 0");
            }
            const amountWei = toWei(amount, "wei");
            const { withdraw } = this.contract.methods;
            withdraw(amountWei).send({ from: this.selectedAccount })
                .once("transactionHash", (txHash) => {
                    this.setStatus(Status.Info, `Transaction ${txHash} sent.`);
                })
                .once("receipt", (receipt) => {
                    this.setStatus(Status.Success, `Transaction complete! Successfully withdrawn ${amountWei} from wallet contract.`);
                })
                .on("error", (err) => {
                    this.setStatus(Status.Fail, `Transaction failed: ${err}`);
                });
        } catch (err) {
            this.setStatus(Status.Fail, `Transaction failed: ${err.message}`);
        }
        await this.refreshPageInfo();
    },

    transfer: async function () {
        const amount = document.getElementById("transferAmount").value;
        const recipient = document.getElementById("transferRecipient").value;

        try {
            if (!this.isValidInputValue(amount)) {
                throw new Error("Amount should be greater than 0");
            }
            if (recipient === null || recipient === "") {
                throw new Error("Recipient address not informed");
            }
            const amountWei = toWei(amount, "wei");
            const { transfer } = this.contract.methods;
            transfer(recipient, amountWei).send({ from: this.selectedAccount })
                .once("transactionHash", (txHash) => {
                    this.setStatus(Status.Info, `Transaction ${txHash} sent.`);
                })
                .once("receipt", (receipt) => {
                    this.setStatus(Status.Success, `Transaction complete! Successfully transferred ${amountWei} wei from contract to ${recipient}.`);
                })
                .on("error", (err) => {
                    this.setStatus(Status.Fail, `Transaction failed: ${err}`);
                });
        } catch (err) {
            this.setStatus(Status.Fail, `Transaction failed: ${err.message}`);
        }
        await this.refreshPageInfo();
    },

    // Helper functions
    refreshWalletBalance: async function () {
        const { getBalance } = this.contract.methods;
        const balance = await getBalance().call();
        const balanceElement = document.getElementById("walletBalance");
        balanceElement.innerHTML = balance;
    },

    refreshAccountsInfo: async function () {
        const accountBalanceElements = document.getElementsByClassName("accountBalance");
        let i = 0;
        for (let account of this.accounts.keys()) {
            let accountBalance = await this.getAccountBalance(account);
            accountBalanceElements[i].innerHTML = accountBalance;
            i++;
        };
    },

    isValidInputValue: function (amount) {
        // number as string ignoring empty string
        return (amount !== null && (/^\s*([0-9]+)\s*$/.test(amount)))
    },

    createCopyButton: function (account) {
        let td = document.createElement("td");
        td.className = "text-end"
        let button = document.createElement("button");
        button.type = "button";
        button.classList = "btn btn-outline-secondary";
        button.onclick = () => {
            navigator.clipboard.writeText(account).then(() => {
                this.setStatus(Status.Info, `address ${account} copied to the clipboard`);
            }, (err) => {
                this.setStatus(Status.Fail, `address copy failed: ${err.message}`);
            });
        };
        let icon = document.createElement("i");
        icon.classList = "fa fa-clipboard";

        button.appendChild(icon);
        td.appendChild(button);
        return td;
    },

    buildAccountList: async function () {
        const accountListElement = document.getElementById("accountList");
        let i = 0;
        let selectedElement = document.getElementById("selectedAccount");
        this.accounts.forEach((balance, account, index) => {
            let tr = document.createElement("tr");
            if (i == 0) {
                tr.classList.add("table-active");
                selectedElement.innerHTML = account;
            }
            let tdAccount = document.createElement("td");
            tdAccount.className = "account";
            tdAccount.innerHTML = account;

            let tdBalance = document.createElement("td");
            tdBalance.className = "accountBalance";
            tdBalance.innerHTML = balance;

            tr.appendChild(tdAccount);
            tr.appendChild(tdBalance);
            // add copy button column
            tr.appendChild(this.createCopyButton(account))
            accountListElement.appendChild(tr);
            i++;
        });

        accountListElement.addEventListener("click", this.selectAccount);
    },

    selectAccount: function (event) {
        let row = event.target.parentNode;
        // if is not a row of the table, ignore
        if (row.matches("tr")) {
            const tbody = document.getElementById("accountList");
            let current = tbody.getElementsByClassName("table-active");
            if (current !== null && current.length > 0) {
                // remove current active state
                current[0].classList.remove("table-active");

                // mark selected account as active
                let selected = row.getElementsByClassName("account").item(0).innerText;
                row.classList.add("table-active");

                // update selected account banner
                let selectedElement = document.getElementById("selectedAccount");
                selectedElement.innerHTML = selected;

                // update default sender account
                this.selectedAccount = selected;
            }
        }
    },

    refreshPageInfo: async function () {
        await this.refreshWalletBalance();
        await this.refreshAccountsInfo();
    },

    setStatus: function (type, message) {
        const status = document.getElementById("notificationPlaceholder");
        let notification = document.createElement("div");
        let closeNotification = function () {
            status.removeChild(notification);
        };
        notification.addEventListener("click", closeNotification, false);
        notification.innerHTML = `
        <div class="alert alert-${type} alert-dismissible" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>`;

        status.append(notification);
        // workaround to close notification popups after 5 seconds
        setTimeout(closeNotification, 5000);
        this.refreshPageInfo();
    },

    getAccounts: async function () {
        let accountsBalances = new Map();
        let accounts = await this.provider.eth.getAccounts();
        for (var i = 0; i < accounts.length; i++) {
            let accountBalance = await this.getAccountBalance(accounts[i]);
            accountsBalances.set(accounts[i], accountBalance)
        }
        return accountsBalances;
    },

    setDefaultAccount: async function () {
        let accounts = await this.provider.eth.getAccounts();
        this.selectedAccount = accounts[0]; // set default sender account

        // retrieve contract owner account
        const { owner } = this.contract.methods;
        this.ownerAccount = await owner().call({ from: this.selectedAccount });
        let ownerElement = document.getElementById("ownerAccount");
        ownerElement.innerHTML = this.ownerAccount;
    },

    getAccountBalance: async function (account) {
        const balance = await this.provider.eth.getBalance(account);
        return fromWei(balance, "ether");
    },
};

window.App = App;

const loadProvider = function () {
    if (window.ethereum) {
        console.log("Using Metamask provider");
        window.ethereum.enable(); // get permission to access accounts
        return new Web3(window.ethereum);
    }
    console.log("No provider detected. Falling back to Web3 at http://127.0.0.1:7545");
    // fallback to local node (ganache)
    return new Web3(new Web3.providers.HttpProvider("http://127.0.0.1:7545"));
};

window.addEventListener("load", function () {
    App.provider = loadProvider();
    App.start();
});
