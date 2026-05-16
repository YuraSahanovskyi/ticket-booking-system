TRUNCATE public.events, public.seats CASCADE;

INSERT INTO public.events (id, title, description, location, start_time, end_time)
VALUES 
    (uuid_generate_v4(), 'Весняний джазовий вечір', 'Жива музика у виконанні кращих джаз-бендів міста.', 'Київ, Будинок Кіно', NOW() + INTERVAL '5 days', NOW() + INTERVAL '5 days 3 hours'),
    (uuid_generate_v4(), 'Виставка сучасного мистецтва', 'Експозиція робіт молодих художників.', 'Галерея М17', NOW() + INTERVAL '8 days', NOW() + INTERVAL '15 days'),
    (uuid_generate_v4(), 'Благодійна вистава "Лісова пісня"', 'Класична українська драма.', 'Театр Лесі Українки', NOW() + INTERVAL '12 days', NOW() + INTERVAL '12 days 2 hours'),
    (uuid_generate_v4(), 'Гастрономічний фестиваль', 'Дегустація страв від провідних шеф-кухарів.', 'ВДНГ', NOW() + INTERVAL '15 days', NOW() + INTERVAL '17 days'),
    (uuid_generate_v4(), 'Кінопоказ просто неба', 'Перегляд класики світового кінематографа.', 'Парк Наталка', NOW() + INTERVAL '18 days', NOW() + INTERVAL '18 days 2 hours'),
    (uuid_generate_v4(), 'Дитяче шоу фокусів', 'Захоплива програма для всієї родини.', 'Цирк', NOW() + INTERVAL '20 days', NOW() + INTERVAL '20 days 1 hour'),
    (uuid_generate_v4(), 'Майстер-клас із гончарства', 'Створення виробів із глини.', 'Арт-завод Платформа', NOW() + INTERVAL '22 days', NOW() + INTERVAL '22 days 4 hours'),
    (uuid_generate_v4(), 'Вечір класичної поезії', 'Читання віршів у затишній атмосфері.', 'Книгарня Є', NOW() + INTERVAL '25 days', NOW() + INTERVAL '25 days 2 hours'),
    (uuid_generate_v4(), 'Спортивний марафон', 'Забіг на 5 та 10 кілометрів.', 'Набережна', NOW() + INTERVAL '30 days', NOW() + INTERVAL '30 days 5 hours'),
    (uuid_generate_v4(), 'Великий стендап-концерт', 'Вечір гумору від відомих коміків.', 'Жовтневий палац', NOW() + INTERVAL '35 days', NOW() + INTERVAL '35 days 2 hours');

INSERT INTO public.seats (id, event_id, "row", number, price)
SELECT 
    uuid_generate_v4(),
    e.id, 
    r.r, 
    n.n, 
    CASE 
        WHEN r.r <= 2 THEN 500
        ELSE 300 
    END
FROM public.events e
CROSS JOIN (SELECT generate_series(1, 6) AS r) r
CROSS JOIN (SELECT generate_series(1, 8) AS n) n;