
create table if not exists projects (
    id integer primary key autoincrement,
    name nvarchar(255) not null,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp
);

-- [SEP] --

create unique index if not exists idx_project_name on projects (name);

-- [SEP] --

create table if not exists versions (
    id integer primary key autoincrement,
    project_id integer not null foreign key references projects(id) on delete cascade,
    digest nvarchar(128) not null,
    digest_type nvarchar(16) not null, -- Either SHA256, MD5, or BLAKE2_256
    filepath nvarchar(256) not null,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp,
)

-- [SEP] --

create table if not exists version_metadata_fields (
    id integer primary key autoincrement,
    version_id integer not null foreign key references versions(id) on delete cascade,
    key nvarchar(255) not null,
    value nvarchar(1024) not null,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp
);
