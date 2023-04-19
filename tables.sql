-- Table Definitions

-- DSBlockMinerInfo
CREATE SEQUENCE IF NOT EXISTS DSBlockMinerInfo_id_seq;
CREATE TABLE "DSBlockMinerInfo" (
    "id" int4 NOT NULL DEFAULT nextval('DSBlockMinerInfo_id_seq'::regclass),
    "dsblock" int,
    "dscommittee" bool,
    "node" bool,
    "walletpublickey" varchar(255),
    "walletaddress" varchar(255),
    "walletaddressbech32" varchar(255),
    PRIMARY KEY ("id")
);

-- Transactions
CREATE SEQUENCE IF NOT EXISTS Tx_id_seq;
CREATE TABLE "Transactions" (
    "id" int4 NOT NULL DEFAULT nextval('Tx_id_seq'::regclass),
    "dsblock" int,
    "txblock" int,
    "txid" varchar(255),
	"amount" float,
	"success" bool,
	"senderpublickey" varchar(255),
	"senderaddressbech32" varchar(255),
	"senderaddress" varchar(255),
	"toaddress" varchar(255),
	"toaddressbech32" varchar(255),
    PRIMARY KEY ("id")
);