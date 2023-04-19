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

Here are some helpful queries to analyse the data/link it up

Blocks I looked at - 28112 DS Block for Rewards and 2811300-2811399 for TX
- To get the mining transactions and wallets that are sent zil from the miners

select t.toaddress, count(*) from public."DSBlockMinerInfo" ds
join public."Transactions" t on ds.walletpublickey = t.senderpublickey
group by t.toaddress

- Transactions in the 100 block that are zil only tx but not all related to mining.

select * from public."DSBlockMinerInfo" ds
right join public."Transactions" t on ds.walletpublickey = t.senderpublickey
where ds.dscommittee is null

- Count of miners in the block

select count (*) from "DSBlockMinerInfo"  where dsblock = xxx

- Average amount for the miners tx

select avg(t.amount) from public."DSBlockMinerInfo" ds
join public."Transactions" t on ds.walletpublickey = t.senderpublickey

- Grouping the tx bewteen 40 and 60 zil that werent part of that blocks rewards

select t.toaddress, count(*) from public."DSBlockMinerInfo" ds
right join public."Transactions" t on ds.walletpublickey = t.senderpublickey
where ds.dscommittee is null and t.amount > 40 and t.amount < 60
group by t.toaddress
order by count(*)

- Excluding the 5 main wallets, looking at the sender of these rather than toaddress

select t.senderaddress, count(*) from public."DSBlockMinerInfo" ds
right join public."Transactions" t on ds.walletpublickey = t.senderpublickey
where ds.dscommittee is null and t.amount > 40 and t.amount < 60 and t.toaddress
not in ('zil1n6yhtv9zrlts8raqhgnr5r2dhyfmyel8egm87c',
zil12chmjxuhs2alj0m2zngu3tjxl2t95zweynx8sl',
zil1ayakns9zz8aemxmwcsjamu6ky97pxh86p4tk70',
zil17rpyuuf3vw4z0jfhqyc04fw2mnpffwj9w9na5p',
zil15xvtse0rvcfwxetstvun72kw5daz8kge0frn3y')
group by t.senderaddress
order by count(*)

- Total NonContract TX in the 100 blocks

select count(*) from "Transactions"

- Total NonContract TX in the 100 blocks not to the 6 main mining addresses

select count(*) from "Transactions"
where senderaddress in ('zil1n6yhtv9zrlts8raqhgnr5r2dhyfmyel8egm87c',
zil12chmjxuhs2alj0m2zngu3tjxl2t95zweynx8sl',
zil1ayakns9zz8aemxmwcsjamu6ky97pxh86p4tk70',
zil17rpyuuf3vw4z0jfhqyc04fw2mnpffwj9w9na5p',
zil15xvtse0rvcfwxetstvun72kw5daz8kge0frn3y',
zil1s5zg376kx586q72dlum497heexkxqdygsr8jpx')
or
toaddress in ('zil1n6yhtv9zrlts8raqhgnr5r2dhyfmyel8egm87c',
zil12chmjxuhs2alj0m2zngu3tjxl2t95zweynx8sl',
zil1ayakns9zz8aemxmwcsjamu6ky97pxh86p4tk70',
zil17rpyuuf3vw4z0jfhqyc04fw2mnpffwj9w9na5p',
zil15xvtse0rvcfwxetstvun72kw5daz8kge0frn3y',
zil1s5zg376kx586q72dlum497heexkxqdygsr8jpx')

- Of what tx are remaining, group by sending address

select senderaddress, count(*) from "Transactions"
where senderaddress not in ('zil1n6yhtv9zrlts8raqhgnr5r2dhyfmyel8egm87c',
zil12chmjxuhs2alj0m2zngu3tjxl2t95zweynx8sl',
zil1ayakns9zz8aemxmwcsjamu6ky97pxh86p4tk70',
zil17rpyuuf3vw4z0jfhqyc04fw2mnpffwj9w9na5p',
zil15xvtse0rvcfwxetstvun72kw5daz8kge0frn3y',
zil1s5zg376kx586q72dlum497heexkxqdygsr8jpx')
and
toaddress not in ('zil1n6yhtv9zrlts8raqhgnr5r2dhyfmyel8egm87c',
zil12chmjxuhs2alj0m2zngu3tjxl2t95zweynx8sl',
zil1ayakns9zz8aemxmwcsjamu6ky97pxh86p4tk70',
zil17rpyuuf3vw4z0jfhqyc04fw2mnpffwj9w9na5p',
zil15xvtse0rvcfwxetstvun72kw5daz8kge0frn3y',
zil1s5zg376kx586q72dlum497heexkxqdygsr8jpx')
group by senderaddress
order by count(*) desc
