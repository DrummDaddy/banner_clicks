CREATE TABLE banner_clicks (
    ts timestamp not null, 
    banner_id int nor null, 
    cnt int not null default 1, 
    primary key (ts, banner_id) 

);