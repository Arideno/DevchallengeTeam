create table countries
(
    id   serial not null
        constraint countries_pkey
            primary key,
    name varchar
        constraint countries_name_key
            unique,
    code varchar(3)
        constraint countries_code_key
            unique
);

alter table countries
    owner to postgres;

INSERT INTO public.countries (id, name, code) VALUES (1, 'Австралія', 'au');
INSERT INTO public.countries (id, name, code) VALUES (2, 'Австрія', 'at');
INSERT INTO public.countries (id, name, code) VALUES (3, 'Азербайджан', 'az');
INSERT INTO public.countries (id, name, code) VALUES (4, 'Алжир', 'dz');
INSERT INTO public.countries (id, name, code) VALUES (5, 'Ангола', 'ao');
INSERT INTO public.countries (id, name, code) VALUES (6, 'Аргентина', 'ar');
INSERT INTO public.countries (id, name, code) VALUES (7, 'Бельгія', 'be');
INSERT INTO public.countries (id, name, code) VALUES (8, 'Білорусь', 'by');
INSERT INTO public.countries (id, name, code) VALUES (9, 'Болгарія', 'bg');
INSERT INTO public.countries (id, name, code) VALUES (10, 'Боснія і Герцеговина', 'ba');
INSERT INTO public.countries (id, name, code) VALUES (11, 'Бразилія', 'br');
INSERT INTO public.countries (id, name, code) VALUES (12, 'Велика Британія', 'gb');
INSERT INTO public.countries (id, name, code) VALUES (13, 'В’єтнам', 'vn');
INSERT INTO public.countries (id, name, code) VALUES (14, 'Вірменія', 'am');
INSERT INTO public.countries (id, name, code) VALUES (15, 'Греція', 'gr');
INSERT INTO public.countries (id, name, code) VALUES (16, 'Грузія', 'ge');
INSERT INTO public.countries (id, name, code) VALUES (17, 'Данія', 'dk');
INSERT INTO public.countries (id, name, code) VALUES (18, 'Естонія', 'ee');
INSERT INTO public.countries (id, name, code) VALUES (19, 'Ефіопія', 'et');
INSERT INTO public.countries (id, name, code) VALUES (20, 'Єгипет', 'eg');
INSERT INTO public.countries (id, name, code) VALUES (21, 'Ізраїль', 'il');
INSERT INTO public.countries (id, name, code) VALUES (22, 'Індонезія', 'id');
INSERT INTO public.countries (id, name, code) VALUES (23, 'Ірак', 'iq');
INSERT INTO public.countries (id, name, code) VALUES (24, 'Іран', 'ir');
INSERT INTO public.countries (id, name, code) VALUES (25, 'Іспанія', 'es');
INSERT INTO public.countries (id, name, code) VALUES (26, 'Італія', 'it');
INSERT INTO public.countries (id, name, code) VALUES (27, 'Йорданія', 'jo');
INSERT INTO public.countries (id, name, code) VALUES (28, 'Казахстан', 'kz');
INSERT INTO public.countries (id, name, code) VALUES (29, 'Канада', 'ca');
INSERT INTO public.countries (id, name, code) VALUES (30, 'Катар', 'qa');
INSERT INTO public.countries (id, name, code) VALUES (31, 'КНР', 'cn');
INSERT INTO public.countries (id, name, code) VALUES (32, 'Кенія', 'ke');
INSERT INTO public.countries (id, name, code) VALUES (33, 'Кіпр', 'cy');
INSERT INTO public.countries (id, name, code) VALUES (34, 'Корея', 'kr');
INSERT INTO public.countries (id, name, code) VALUES (35, 'Куба', 'cu');
INSERT INTO public.countries (id, name, code) VALUES (36, 'Кувейт', 'kw');
INSERT INTO public.countries (id, name, code) VALUES (37, 'Латвія', 'lv');
INSERT INTO public.countries (id, name, code) VALUES (38, 'Литва', 'lt');
INSERT INTO public.countries (id, name, code) VALUES (39, 'Ліван', 'lb');
INSERT INTO public.countries (id, name, code) VALUES (40, 'Малайзія', 'my');
INSERT INTO public.countries (id, name, code) VALUES (41, 'Північна Македонія', 'mk');
INSERT INTO public.countries (id, name, code) VALUES (42, 'Марокко', 'ma');
INSERT INTO public.countries (id, name, code) VALUES (43, 'Мексика', 'mx');
INSERT INTO public.countries (id, name, code) VALUES (44, 'Молдова', 'md');
INSERT INTO public.countries (id, name, code) VALUES (45, 'Нігерія', 'ng');
INSERT INTO public.countries (id, name, code) VALUES (46, 'Нідерланди', 'nl');
INSERT INTO public.countries (id, name, code) VALUES (47, 'Німеччина', 'de');
INSERT INTO public.countries (id, name, code) VALUES (48, 'Норвегія', 'no');
INSERT INTO public.countries (id, name, code) VALUES (49, 'ОАЕ', 'ae');
INSERT INTO public.countries (id, name, code) VALUES (50, 'Пакистан', 'pk');
INSERT INTO public.countries (id, name, code) VALUES (51, 'ПАР', 'za');
INSERT INTO public.countries (id, name, code) VALUES (52, 'Перу', 'pe');
INSERT INTO public.countries (id, name, code) VALUES (53, 'Польща', 'pl');
INSERT INTO public.countries (id, name, code) VALUES (54, 'Португалія', 'pt');
INSERT INTO public.countries (id, name, code) VALUES (55, 'РФ', 'ru');
INSERT INTO public.countries (id, name, code) VALUES (56, 'Румунія', 'ro');
INSERT INTO public.countries (id, name, code) VALUES (57, 'Саудівська Аравія', 'sa');
INSERT INTO public.countries (id, name, code) VALUES (58, 'Сенегал', 'sn');
INSERT INTO public.countries (id, name, code) VALUES (59, 'Сінгапур', 'sg');
INSERT INTO public.countries (id, name, code) VALUES (60, 'Словакія', 'sk');
INSERT INTO public.countries (id, name, code) VALUES (61, 'Словенія', 'si');
INSERT INTO public.countries (id, name, code) VALUES (62, 'США', 'us');
INSERT INTO public.countries (id, name, code) VALUES (63, 'Таджикистан', 'tj');
INSERT INTO public.countries (id, name, code) VALUES (64, 'Таїланд', 'th');
INSERT INTO public.countries (id, name, code) VALUES (65, 'Туніс', 'tn');
INSERT INTO public.countries (id, name, code) VALUES (66, 'Лівія', 'ly');
INSERT INTO public.countries (id, name, code) VALUES (67, 'Туреччина', 'tr');
INSERT INTO public.countries (id, name, code) VALUES (68, 'Туркменістан', 'tm');
INSERT INTO public.countries (id, name, code) VALUES (69, 'Угорщина', 'hu');
INSERT INTO public.countries (id, name, code) VALUES (70, 'Узбекистан', 'uz');
INSERT INTO public.countries (id, name, code) VALUES (71, 'Чорногорія', 'me');
INSERT INTO public.countries (id, name, code) VALUES (72, 'Фінляндія', 'fi');
INSERT INTO public.countries (id, name, code) VALUES (73, 'Франція', 'fr');
INSERT INTO public.countries (id, name, code) VALUES (74, 'Хорватія', 'hr');
INSERT INTO public.countries (id, name, code) VALUES (75, 'Чилі', 'cl');
INSERT INTO public.countries (id, name, code) VALUES (76, 'Швеція', 'se');
INSERT INTO public.countries (id, name, code) VALUES (77, 'Японія', 'jp');