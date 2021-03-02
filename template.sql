create table billings
(
    id                 int auto_increment
        primary key,
    userId             int         default 0 null,
    transactionId      varchar(50) default '' null,
    productId          varchar(50) default '' null,
    purchaseToken      text null,
    purchaseDate       datetime null,
    created            datetime    default current_timestamp() null,
    updated            datetime    default current_timestamp() null on update current_timestamp (),
    deleted            datetime null,
    purchaseState      tinyint(1) default 0 null,
    useState           tinyint     default 0 null,
    refundRequestState tinyint     default 0 null,
    allowState         tinyint(1) default 1 null
);

create
index userId
    on billings (userId);

create table blocks
(
    id          int(11) unsigned auto_increment
        primary key,
    login_type  tinyint(1) unsigned default 0 null,
    uid         varchar(21)  default '' null,
    uuid        varchar(40)  default '' null,
    name        varchar(8)   default '' null,
    description varchar(255) default '' null,
    date        datetime     default current_timestamp() null
);

create table items
(
    id          int(11) unsigned auto_increment
        primary key,
    name        varchar(50) default '' null,
    description varchar(50) default '' null,
    icon        varchar(50) default '' null,
    type        tinyint unsigned default 0 null,
    num         tinyint unsigned default 0 null,
    `range`     tinyint unsigned default 0 null,
    speed       tinyint unsigned default 0 null,
    cool        tinyint unsigned default 0 null,
    cost        int         default 0 null,
    method      varchar(50) default '' null
);

create table portals
(
    no         int(11) unsigned auto_increment,
    place      smallint(3) unsigned null,
    x          smallint(3) unsigned null,
    y          smallint(3) unsigned null,
    next_place smallint(3) unsigned null,
    next_x     smallint(3) unsigned null,
    next_y     smallint(3) unsigned null,
    next_dir_x tinyint(1) unsigned zerofill default 0 not null,
    next_dir_y tinyint(1) unsigned zerofill default 0 not null,
    sound      varchar(50) default '' not null,
    constraint no_UNIQUE
        unique (no)
);

create
index no
    on portals (no);

alter table portals
    add primary key (no);

create table users
(
    id            int(11) unsigned auto_increment
        primary key,
    login_type    tinyint(1) unsigned null,
    uid           varchar(21) null,
    uuid          varchar(40) null,
    email         varchar(150) null,
    name          varchar(8) null,
    level         smallint(3) unsigned default 1 null,
    exp           int(11) unsigned default 0 null,
    coin          int(11) unsigned default 1000 null,
    cash          int(11) unsigned default 0 null,
    win           int(11) unsigned default 0 null,
    lose          int(11) unsigned default 0 null,
    escape        int(11) unsigned default 0 null,
    `kill`        int(11) unsigned default 0 null,
    death         int(11) unsigned default 0 null,
    assist        int(11) unsigned default 0 null,
    blast         int(11) unsigned default 0 null,
    rescue        int(11) unsigned default 0 null,
    rescue_combo  int(11) unsigned default 0 null,
    survive       int(11) unsigned default 0 null,
    red_graphics  varchar(50)  default 'ao' null,
    blue_graphics varchar(50)  default 'Mania' null,
    memo          varchar(255) default '' null,
    last_chat     datetime     default current_timestamp() null,
    dickhead      tinyint(1) unsigned default 0 null,
    permission    tinyint(1) unsigned default 1 null,
    admin         tinyint(1) unsigned default 0 null,
    verify        tinyint(1) default 0 null,
    sex           int null,
    point         int null,
    constraint nickname
        unique (name)
);

create table clans
(
    id        int(11) unsigned auto_increment
        primary key,
    master_id int(11) unsigned default 0 not null,
    name      varchar(50) default ''                  not null,
    level     tinyint unsigned default 1 not null,
    exp       int(11) unsigned default 0 not null,
    coin      int(11) unsigned default 0 not null,
    regdate   datetime    default current_timestamp() not null,
    constraint name
        unique (name),
    constraint FK_clans_users
        foreign key (master_id) references users (id)
);

create table clan_members
(
    id      int(11) unsigned auto_increment
        primary key,
    clan_id int(11) unsigned default 0 null,
    user_id int(11) unsigned default 0 null,
    grade   tinyint(3) default 1 null,
    level   int null,
    exp     int(11) unsigned default 0 null,
    coin    int(11) unsigned default 0 null,
    regdate datetime default current_timestamp() null,
    constraint FK_clan_members_clans
        foreign key (clan_id) references clans (id),
    constraint FK_clan_members_users
        foreign key (user_id) references users (id)
);

create table invite_clans
(
    id        int(11) unsigned auto_increment
        primary key,
    clan_id   int(11) unsigned not null,
    user_id   int(11) unsigned not null,
    target_id int(11) unsigned not null,
    regdate   datetime default current_timestamp() not null,
    constraint FK_invite_clans_clans
        foreign key (clan_id) references clans (id),
    constraint FK_invite_clans_users
        foreign key (user_id) references users (id),
    constraint FK_invite_clans_users_2
        foreign key (target_id) references users (id)
);

create
index uid_login_type
    on users (uid, login_type);

create
index verify
    on users (verify);

