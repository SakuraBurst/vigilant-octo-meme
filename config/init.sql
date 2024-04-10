create table feature (id int primary key);
create table tags (id int primary key);
create table banners_tags (id serial primary key, tag_id int, banner_id int, foreign key (tag_id) references tags (id), foreign key (banner_id) references banners (id));
create table banners (id serial primary key, feature_id int, content json, is_active bool, foreign key (feature_id) references feature (id));


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