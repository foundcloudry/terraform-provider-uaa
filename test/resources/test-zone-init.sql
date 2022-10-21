create extension if not exists "uuid-ossp";

insert into identity_zone (id, name, subdomain)
    values ('test-zone', 'Test Zone', 'test-zone')
    on conflict do nothing;

insert into identity_provider (id, identity_zone_id, name, origin_key, type, config)
select uuid_generate_v4(), 'test-zone', name, origin_key, type, config
from identity_provider
where identity_zone_id='uaa' and name='uaa' and origin_key='uaa' and type='uaa'
    on conflict do nothing;

insert into groups (id, displayname, identity_zone_id, description)
select uuid_generate_v4(), displayname, 'test-zone', description
from groups where identity_zone_id='uaa'
    on conflict do nothing;

