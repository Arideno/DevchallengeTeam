create table topics
(
    id   serial  not null
        constraint topics_pk
            primary key,
    name varchar not null
);

alter table topics
    owner to postgres;

create unique index topics_name_uindex
    on topics (name);

INSERT INTO public.topics (id, name) VALUES (1, 'Консульський прийом');
INSERT INTO public.topics (id, name) VALUES (2, 'Паспортні документи');
INSERT INTO public.topics (id, name) VALUES (3, 'Віза');
INSERT INTO public.topics (id, name) VALUES (4, 'Витребування документів з України');
INSERT INTO public.topics (id, name) VALUES (5, 'Консульський облік');
INSERT INTO public.topics (id, name) VALUES (6, 'Постійне проживання  за кордоном');
INSERT INTO public.topics (id, name) VALUES (7, 'Громадянство');
INSERT INTO public.topics (id, name) VALUES (8, 'Захист громадян України');
INSERT INTO public.topics (id, name) VALUES (9, 'Надзвичайні події');
INSERT INTO public.topics (id, name) VALUES (10, 'Виборчий процес за кордоном');
INSERT INTO public.topics (id, name) VALUES (11, 'Сімейні питання');
INSERT INTO public.topics (id, name) VALUES (12, 'Почесні консули України');
INSERT INTO public.topics (id, name) VALUES (13, 'Легалізація');
INSERT INTO public.topics (id, name) VALUES (14, 'Нотаріальні дії, видача довідок');
INSERT INTO public.topics (id, name) VALUES (15, 'Звернення громадян');