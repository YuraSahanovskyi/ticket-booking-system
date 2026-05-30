TRUNCATE public.events, public.seats CASCADE;

WITH new_event AS (
    INSERT INTO public.events (
        id,
        title,
        description,
        location,
        start_time,
        end_time
    )
    VALUES (
        uuid_generate_v4(),
        'Тестова подія',
        'Подія для тестування',
        'Київ',
        NOW() + INTERVAL '5 days',
        NOW() + INTERVAL '5 days 2 hours'
    )
    RETURNING id
)
INSERT INTO public.seats (
    id,
    event_id,
    "row",
    number,
    price
)
SELECT
    uuid_generate_v4(),
    e.id,
    1,
    s.number,
    300
FROM new_event e
CROSS JOIN generate_series(1, 3) AS s(number);