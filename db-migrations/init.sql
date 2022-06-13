drop table if exists trades;

create table trades
(
    evt timestamp,
    pair varchar(100),
    tid bigint,
    p double precision,
    q double precision,
    b bigint,
    a bigint,
    tt timestamp
)

--todo add ticker table