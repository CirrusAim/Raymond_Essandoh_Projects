module.exports = {
  // Uncommenting the defaults below 
  // provides for an easier quick-start with Ganache.
  // You can also follow this format for other networks;
  // see <http://truffleframework.com/docs/advanced/configuration>
  // for more details on how to specify configuration options!
  networks: {
    development: { // local test net
      host: "127.0.0.1",
      port: 8545,
      network_id: "*"
    },
    ganache: { // ganache-cli
      host: "127.0.0.1",
      port: 7545,
      network_id: "5777"
    },
    develop: { // truffle development
      host: "127.0.0.1",
      port: 8545,
      network_id: "*",
      accounts: 5,
      defaultEtherBalance: 50
    },
  },
  compilers: {
    solc: {
      version: '0.8.0'
      // settings: {
      //  optimizer: {
      //    enabled: false,
      //    runs: 200
      //  },
      //  evmVersion: "byzantium"
      // }
    }
  }
};