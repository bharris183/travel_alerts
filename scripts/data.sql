use threat_alerts;

insert into RSS_THREATS 
    (
        `COUNTRY_CODE` , 
        `THREAT_LEVEL` ,
        `TITLE` ,
        `LINK` ,
        `DESCRIPTION` ,
        `PUB_DATE`
    )
    values
    (
        'SN' ,
        3 ,
        'Senegal is safe, mostly' ,
        'http://sn.com' ,
        'Senegal is a lovely place. Just stick to your hotel.' ,
        '2021-05-11'
    )
