# Sample project to pull mining tx from
1) Interract with the Zilliqa Blockchain API
 - Uses GetTxnBodiesForTxBlock and GetMinerInfo - Reference 
 https://dev.zilliqa.com/api/blockchain-related-methods/api-blockchain-get-miner-info/
 https://dev.zilliqa.com/api/transaction-related-methods/api-transaction-get-txbodies-for-txblock/
2) Save results to a Postgres SQL database 


To use the project locally, you will need go installed.
https://golang.org/doc/install

Also will need a postgres database.

Also a .env file with the following;

DATABASE_URL=postgresql://(postgres connect string)

Then in main.go set the DS Block Number you want to pull.  Your database will need the tables in tables.sql.