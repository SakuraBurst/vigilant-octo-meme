create table if not exists feature (id int primary key);
create table if not exists tags (id int primary key);
create table if not exists banners (id serial primary key, tags int array, feature_id int not null, content json not null, is_active bool not null, created_at timestamp not null default now(), updated_at timestamp not null default now(), foreign key (feature_id) references feature (id));
create table if not exists banners_tags (id serial primary key, tag_id int not null, banner_id int not null, foreign key (tag_id) references tags (id), foreign key (banner_id) references banners (id));

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE or replace TRIGGER set_timestamp
    BEFORE UPDATE ON banners
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

--тегайди 2 феатуреайди 3
--select (content) from banners b
--join banners_tags bt on b.id = bt.banner_id
--where b.feature_id = 3 and bt.tag_id = 2

--тегайди 2
--select (content) from banners b
--join banners_tags bt on b.id = bt.banner_id
--where bt.tag_id = 2

--феатуреайди 3
--select (content) from banners where b.feature_id = 3