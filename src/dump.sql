create table if not exists users
(
	id serial not null,
	username text not null,
	password text not null,
	email text not null
)
;

create unique index users_id_uindex
	on users (id)
;

create unique index users_username_uindex
	on users (username)
;

create unique index users_email_uindex
	on users (email)
;

create table if not exists posts
(
	id serial not null,
	author integer not null
		constraint posts_users_id_fk
			references users (id) on delete cascade,
	content text,
	date timestamp with time zone default now() not null,
	likes integer default 0
)
;

create unique index posts_id_uindex
	on posts (id)
;

create table if not exists posts_likes
(
	id serial not null,
	user_id integer not null
		constraint posts_likes_users_id_fk
			references users (id)
				on delete cascade,
	post_id integer not null
		constraint posts_likes_posts_id_fk
			references posts (id)
				on delete cascade
)
;

create unique index posts_likes_id_uindex
	on posts_likes (id)
;

create table if not exists comments
(
	id serial not null,
	content text not null,
	author integer
		constraint comments_users_id_fk
			references users (id)
				on delete cascade,
	date timestamp with time zone default now() not null,
	likes integer default 0 not null,
	post_id integer not null
		constraint comments_posts_id_fk
			references posts (id)
				on delete cascade
)
;

create unique index comments_id_uindex
	on comments (id)
;

create table if not exists comments_likes
(
	id serial not null,
	comment_id integer not null
		constraint comments_likes_comments_id_fk
			references comments (id)
				on delete cascade,
	user_id integer not null
		constraint comments_likes_users_id_fk
			references users (id)
				on delete cascade
)
;

create unique index comments_likes_id_uindex
	on comments_likes (id)
;

create table if not exists followers
(
	id serial not null
		constraint followers_pkey
			primary key,
	user_followed integer not null
		constraint followers_users_followed_fk
			references users (id)
				on delete cascade,
	user_following integer not null
		constraint followers_users_following_fk
			references users (id)
				on delete cascade
)
;

create unique index followers_id_uindex
	on followers (id)
;

