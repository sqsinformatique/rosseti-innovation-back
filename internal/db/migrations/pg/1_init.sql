-- +goose Up
CREATE TABLE IF NOT EXISTS production.users (
    id serial PRIMARY KEY,
    user_hash character varying(255) DEFAULT '',
    user_email character varying(255) DEFAULT '',
    user_phone character varying(255) DEFAULT '',
    user_role character varying(255) DEFAULT '',
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone,
    CONSTRAINT user_email_unique UNIQUE (user_email),
    CONSTRAINT user_phone_unique UNIQUE (user_phone)
);

INSERT INTO production.users (user_hash, user_email, user_phone, user_role) 
VALUES ('15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225', 'test1@rosseti.ru', '+79169999999', 'USER');
INSERT INTO production.users (user_hash, user_email, user_phone, user_role) 
VALUES ('15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225', 'test2@rosseti.ru', '+79169999998', 'MODERATOR');
INSERT INTO production.users (user_hash, user_email, user_phone, user_role) 
VALUES ('15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225', 'test3@rosseti.ru', '+79169999997', 'EXPERT');
INSERT INTO production.users (user_hash, user_email, user_phone, user_role) 
VALUES ('15e2b0d3c33891ebb0f1ef609ec419420c20e320ce94c65fbc8c3312448eb225', 'test4@rosseti.ru', '+79169999996', 'ADMIN');

CREATE TABLE IF NOT EXISTS production.profiles (
    id INTEGER PRIMARY KEY,
    user_first_name character varying(255) DEFAULT '',
    user_middle_name character varying(255) DEFAULT '',
    user_last_name character varying(255) DEFAULT '',
    user_position character varying(255) DEFAULT '',
    user_company character varying(255) DEFAULT '',
    user_electro_group character varying(255) DEFAULT 'Гр. I',
    user_private_key character varying(2048) DEFAULT '',
    user_public_key character varying(2048) DEFAULT '',
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

INSERT INTO production.profiles (id, user_first_name, user_middle_name, user_last_name, user_position, user_company, user_electro_group, user_private_key, user_public_key) 
VALUES (1, 'Василий', 'Васильевич', 'Васильев', 'Мастер', 'ПАО МОЭСК', 'Гр. V',
'MIIEpAIBAAKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQABAoIBAQCDmons6NJpd9FDToEAU4mZFiGQY4mXv+vfp7w4D2nY4JF+R7+/Y5RNtqlxH2CTyQePpCWQAVw6r5mmzHPi2nDbcPfwWzQxCbP0OpUcxmS2zrQssNRpu1LanbS/buTDA2PWOqsQ7/JaO93+bgXHUje7XQ1wRRY0Byy/UtmSjrxApAqMiFGM4koIu9aOl6CvSb2bxO4MwitV9iS1lMFl4SlaYf9zPtIOJYUQWziRNz1ZK0CKCVnashspM+eNgVw+Z6P5OTOfOPDWBL0ORMvRXLQtoYND8tO9oCS+WV7MaeOk1OP8mP8okrb1o9x9ZGhWIQx7Ue7iaYvzgrUuj0gRtLABAoGBAN+PrpHnaY5DIgGqwBBTzNCaXX1W7n0wALU8qin5jo0cScSttMcM9OLM86BRMSCqYfOPDM0Lak3ZfgybxPCivhwwapsJ7ScP43/BzIFhzIZ2Xtg4P4KrF0OJ/ggfOUV6wgv0OKTzgTQ855lZy3MSKmVb7folXf9obYk+93PzZChvAoGBAMn80SrhuHZF81THBZfbcpJxpKNgTQerX4ia6e/CTT4hjlza8kRGawnWfQOFdmCT3mzdb6tQ+KPsxQWCRtH5DWg84GH/6+lgIUMEB+8OjUoDq6y+XmMWlCePVH1+7swHe7LRlaonwagA7SVn5t8+r5sZc+sIpH1wsc35DOZ+GlIZAoGAd3tBH3WAcqnqeN2bPJ6s7igyIxTc7UdEeZhskXZw+3XM7zKvVVrVXomPA3WhPgYRx6wCeWvKasT8mxx9SuaPmF0//JB3kNLrEZKwC84LEyocUo7tUpbCHjSX8htN7pZHM0BZLb9+pD6QwOK+20cwJW/WZkSmUiSrthhTBENmmj0CgYEAlhIUhju2hYlrRO2ppi4RbeSpYglGshANpr0SWmSOZz8fOrYhkcCP/nsx3s/mJ9M1SsUrFqnOUlyz9WfZnl/gKjYwsB8o8/fMPrJcAq1ZJEid4HaAQjagVNQU/ji0yzo0GaPGAuoO4/fsOgJ8chls91tt2I5PSDPWpyYHA6llfOECgYAhU0FLx2eHwQpFkQyGjgAvnRIzDwB6lKtFwO2BULwu8RBwYgQg1cAr8BCESJkpAHYC4FXEN1PT4qXb3fohGWOoIDWB6LliIUiGEFW5/1rFdLVUnXjA0NCqPQngBE/pW8/tdRnfqkjL07BRbz0RyLa5Me/AZfURkl21flBhvvs7bA==',
'MIIBCgKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQAB'
);
INSERT INTO production.profiles (id, user_first_name, user_middle_name, user_last_name, user_position, user_company, user_electro_group, user_private_key, user_public_key) 
VALUES (2, 'Пётр', 'Петрович', 'Петров', 'Электромонтёр', 'ПАО МОЭСК', 'Гр. IV',
'MIIEpAIBAAKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQABAoIBAQCDmons6NJpd9FDToEAU4mZFiGQY4mXv+vfp7w4D2nY4JF+R7+/Y5RNtqlxH2CTyQePpCWQAVw6r5mmzHPi2nDbcPfwWzQxCbP0OpUcxmS2zrQssNRpu1LanbS/buTDA2PWOqsQ7/JaO93+bgXHUje7XQ1wRRY0Byy/UtmSjrxApAqMiFGM4koIu9aOl6CvSb2bxO4MwitV9iS1lMFl4SlaYf9zPtIOJYUQWziRNz1ZK0CKCVnashspM+eNgVw+Z6P5OTOfOPDWBL0ORMvRXLQtoYND8tO9oCS+WV7MaeOk1OP8mP8okrb1o9x9ZGhWIQx7Ue7iaYvzgrUuj0gRtLABAoGBAN+PrpHnaY5DIgGqwBBTzNCaXX1W7n0wALU8qin5jo0cScSttMcM9OLM86BRMSCqYfOPDM0Lak3ZfgybxPCivhwwapsJ7ScP43/BzIFhzIZ2Xtg4P4KrF0OJ/ggfOUV6wgv0OKTzgTQ855lZy3MSKmVb7folXf9obYk+93PzZChvAoGBAMn80SrhuHZF81THBZfbcpJxpKNgTQerX4ia6e/CTT4hjlza8kRGawnWfQOFdmCT3mzdb6tQ+KPsxQWCRtH5DWg84GH/6+lgIUMEB+8OjUoDq6y+XmMWlCePVH1+7swHe7LRlaonwagA7SVn5t8+r5sZc+sIpH1wsc35DOZ+GlIZAoGAd3tBH3WAcqnqeN2bPJ6s7igyIxTc7UdEeZhskXZw+3XM7zKvVVrVXomPA3WhPgYRx6wCeWvKasT8mxx9SuaPmF0//JB3kNLrEZKwC84LEyocUo7tUpbCHjSX8htN7pZHM0BZLb9+pD6QwOK+20cwJW/WZkSmUiSrthhTBENmmj0CgYEAlhIUhju2hYlrRO2ppi4RbeSpYglGshANpr0SWmSOZz8fOrYhkcCP/nsx3s/mJ9M1SsUrFqnOUlyz9WfZnl/gKjYwsB8o8/fMPrJcAq1ZJEid4HaAQjagVNQU/ji0yzo0GaPGAuoO4/fsOgJ8chls91tt2I5PSDPWpyYHA6llfOECgYAhU0FLx2eHwQpFkQyGjgAvnRIzDwB6lKtFwO2BULwu8RBwYgQg1cAr8BCESJkpAHYC4FXEN1PT4qXb3fohGWOoIDWB6LliIUiGEFW5/1rFdLVUnXjA0NCqPQngBE/pW8/tdRnfqkjL07BRbz0RyLa5Me/AZfURkl21flBhvvs7bA==',
'MIIBCgKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQAB'
);
INSERT INTO production.profiles (id, user_first_name, user_middle_name, user_last_name, user_position, user_company, user_electro_group, user_private_key, user_public_key) 
VALUES (3, 'Геннадий', 'Генадьевич', 'Геннадьев', 'Главный инженер', 'ПАО МОЭСК', 'Гр. V',
'MIIEpAIBAAKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQABAoIBAQCDmons6NJpd9FDToEAU4mZFiGQY4mXv+vfp7w4D2nY4JF+R7+/Y5RNtqlxH2CTyQePpCWQAVw6r5mmzHPi2nDbcPfwWzQxCbP0OpUcxmS2zrQssNRpu1LanbS/buTDA2PWOqsQ7/JaO93+bgXHUje7XQ1wRRY0Byy/UtmSjrxApAqMiFGM4koIu9aOl6CvSb2bxO4MwitV9iS1lMFl4SlaYf9zPtIOJYUQWziRNz1ZK0CKCVnashspM+eNgVw+Z6P5OTOfOPDWBL0ORMvRXLQtoYND8tO9oCS+WV7MaeOk1OP8mP8okrb1o9x9ZGhWIQx7Ue7iaYvzgrUuj0gRtLABAoGBAN+PrpHnaY5DIgGqwBBTzNCaXX1W7n0wALU8qin5jo0cScSttMcM9OLM86BRMSCqYfOPDM0Lak3ZfgybxPCivhwwapsJ7ScP43/BzIFhzIZ2Xtg4P4KrF0OJ/ggfOUV6wgv0OKTzgTQ855lZy3MSKmVb7folXf9obYk+93PzZChvAoGBAMn80SrhuHZF81THBZfbcpJxpKNgTQerX4ia6e/CTT4hjlza8kRGawnWfQOFdmCT3mzdb6tQ+KPsxQWCRtH5DWg84GH/6+lgIUMEB+8OjUoDq6y+XmMWlCePVH1+7swHe7LRlaonwagA7SVn5t8+r5sZc+sIpH1wsc35DOZ+GlIZAoGAd3tBH3WAcqnqeN2bPJ6s7igyIxTc7UdEeZhskXZw+3XM7zKvVVrVXomPA3WhPgYRx6wCeWvKasT8mxx9SuaPmF0//JB3kNLrEZKwC84LEyocUo7tUpbCHjSX8htN7pZHM0BZLb9+pD6QwOK+20cwJW/WZkSmUiSrthhTBENmmj0CgYEAlhIUhju2hYlrRO2ppi4RbeSpYglGshANpr0SWmSOZz8fOrYhkcCP/nsx3s/mJ9M1SsUrFqnOUlyz9WfZnl/gKjYwsB8o8/fMPrJcAq1ZJEid4HaAQjagVNQU/ji0yzo0GaPGAuoO4/fsOgJ8chls91tt2I5PSDPWpyYHA6llfOECgYAhU0FLx2eHwQpFkQyGjgAvnRIzDwB6lKtFwO2BULwu8RBwYgQg1cAr8BCESJkpAHYC4FXEN1PT4qXb3fohGWOoIDWB6LliIUiGEFW5/1rFdLVUnXjA0NCqPQngBE/pW8/tdRnfqkjL07BRbz0RyLa5Me/AZfURkl21flBhvvs7bA==',
'MIIBCgKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQAB'
);
INSERT INTO production.profiles (id, user_first_name, user_middle_name, user_last_name, user_position, user_company, user_electro_group, user_private_key, user_public_key) 
VALUES (4, 'Иван', 'Иванович', 'Иванов', 'Администратор', 'ПАО МОЭСК', 'Гр. III',
'MIIEpAIBAAKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQABAoIBAQCDmons6NJpd9FDToEAU4mZFiGQY4mXv+vfp7w4D2nY4JF+R7+/Y5RNtqlxH2CTyQePpCWQAVw6r5mmzHPi2nDbcPfwWzQxCbP0OpUcxmS2zrQssNRpu1LanbS/buTDA2PWOqsQ7/JaO93+bgXHUje7XQ1wRRY0Byy/UtmSjrxApAqMiFGM4koIu9aOl6CvSb2bxO4MwitV9iS1lMFl4SlaYf9zPtIOJYUQWziRNz1ZK0CKCVnashspM+eNgVw+Z6P5OTOfOPDWBL0ORMvRXLQtoYND8tO9oCS+WV7MaeOk1OP8mP8okrb1o9x9ZGhWIQx7Ue7iaYvzgrUuj0gRtLABAoGBAN+PrpHnaY5DIgGqwBBTzNCaXX1W7n0wALU8qin5jo0cScSttMcM9OLM86BRMSCqYfOPDM0Lak3ZfgybxPCivhwwapsJ7ScP43/BzIFhzIZ2Xtg4P4KrF0OJ/ggfOUV6wgv0OKTzgTQ855lZy3MSKmVb7folXf9obYk+93PzZChvAoGBAMn80SrhuHZF81THBZfbcpJxpKNgTQerX4ia6e/CTT4hjlza8kRGawnWfQOFdmCT3mzdb6tQ+KPsxQWCRtH5DWg84GH/6+lgIUMEB+8OjUoDq6y+XmMWlCePVH1+7swHe7LRlaonwagA7SVn5t8+r5sZc+sIpH1wsc35DOZ+GlIZAoGAd3tBH3WAcqnqeN2bPJ6s7igyIxTc7UdEeZhskXZw+3XM7zKvVVrVXomPA3WhPgYRx6wCeWvKasT8mxx9SuaPmF0//JB3kNLrEZKwC84LEyocUo7tUpbCHjSX8htN7pZHM0BZLb9+pD6QwOK+20cwJW/WZkSmUiSrthhTBENmmj0CgYEAlhIUhju2hYlrRO2ppi4RbeSpYglGshANpr0SWmSOZz8fOrYhkcCP/nsx3s/mJ9M1SsUrFqnOUlyz9WfZnl/gKjYwsB8o8/fMPrJcAq1ZJEid4HaAQjagVNQU/ji0yzo0GaPGAuoO4/fsOgJ8chls91tt2I5PSDPWpyYHA6llfOECgYAhU0FLx2eHwQpFkQyGjgAvnRIzDwB6lKtFwO2BULwu8RBwYgQg1cAr8BCESJkpAHYC4FXEN1PT4qXb3fohGWOoIDWB6LliIUiGEFW5/1rFdLVUnXjA0NCqPQngBE/pW8/tdRnfqkjL07BRbz0RyLa5Me/AZfURkl21flBhvvs7bA==',
'MIIBCgKCAQEAsGSYKiaHzXTqMcIIdGE5iDU/3dlEVNMeUogwFSu83DN9JVNsmlIyhRVmEnm55sB6bs4ES0Lph28EeSmwf71sLWmmtWCxp90Cq3gBzDTzQJqyPNfL3eeFeeNDgNtUs3osB4gMdWs+CvbphRbOsAq3sbNQJ00aGZ6p6zkX/IZxrBIjxY+kxWIjQLy0p/Yn4rybM4+VFKMEUCiXM7Deu16N5hz1FF8HtSxsQBc4gszGWvIaVhZ1iK3y7DRqQBpjHCaA6cKC3Nh6ZpTxyIUphKrusBowYnXP5R/e3FQDzh5EhIwzIvEijJx9ywoeF9EC36Tnz62CRUIPbfla09oVlwGA1wIDAQAB'
);

CREATE TABLE IF NOT EXISTS production.sessions (
    id character varying(255) PRIMARY KEY,
    user_id INTEGER NOT NULL,
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

-- +goose Down
DROP TABLE production.users;
DROP TABLE production.profiles;
DROP TABLE production.sessions;
