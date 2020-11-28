-- +goose Up
CREATE TABLE IF NOT EXISTS production.direction (
    id serial PRIMARY KEY,
    title character varying(255) DEFAULT '',
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

INSERT INTO production.direction (title) 
VALUES ('Эксплуатация подстанций (подстанционного оборудования)');
INSERT INTO production.direction (title)
VALUES ('Эксплуатация магистральных сетей');
INSERT INTO production.direction (title)
VALUES ('Эксплуатация распределительных сетей');
INSERT INTO production.direction (title)
VALUES ('Капитальное строительство, реконструкция, проектирование');
INSERT INTO production.direction (title)
VALUES ('Эксплуатация зданий, сооружений, специальной техники');
INSERT INTO production.direction (title)
VALUES ('Оперативно-диспетчерское управление');
INSERT INTO production.direction (title)
VALUES ('Релейная защита и противоаварийная автоматика');
INSERT INTO production.direction (title)
VALUES ('Информационные технологии, системы связи');
INSERT INTO production.direction (title)
VALUES ('Мониторинг и диагностика');
INSERT INTO production.direction (title)
VALUES ('Контроль качества и учёт электроэнергии');
INSERT INTO production.direction (title)
VALUES ('Производственная безопасность и охрана труда');
INSERT INTO production.direction (title)
VALUES ('Технологическое присоединение');
INSERT INTO production.direction (title)
VALUES ('Аварийно-восстановительные работы');
INSERT INTO production.direction (title)
VALUES ('Экология, энергоэффективность, снижение потерь');
INSERT INTO production.direction (title)
VALUES ('Совершенствование системы управления');
INSERT INTO production.direction (title)
VALUES ('Дополнительные (нетарифные) услуги');

CREATE TABLE IF NOT EXISTS production.theme (
    id serial PRIMARY KEY,
    author_id INTEGER NOT NULL,
    title character varying(255) DEFAULT '',
    direction INTEGER NOT NULL,
    tags character varying(255) DEFAULT '',
    like_counter INTEGER NOT NULL DEFAULT 0,
    meta jsonb,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    deleted_at timestamp with time zone
);

-- +goose Down
DROP TABLE production.direction;
DROP TABLE production.theme;
