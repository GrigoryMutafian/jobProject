CREATE TABLE IF NOT EXISTS subs_table (
    id SERIAL PRIMARY KEY,
    service VARCHAR(50) NOT NULL,
    Price INT DEFAULT 0,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE
);

INSERT INTO subs_table (service, price, user_id, start_date, end_date)
VALUES
  ('Yandex', 299, '70601fee-2bf1-4721-ae6f-7636e79a0cbb', TO_DATE('10-2025','MM-YYYY'), TO_DATE('11-2025','MM-YYYY')),
  ('Кинопоиск', 399, '60601fee-2bf1-4721-ae6f-7636e79a0cba', TO_DATE('10-2025','MM-YYYY'), NULL),
  ('IVI', 499, '90601fee-32f1-4721-ae6f-7636e79a0cba', TO_DATE('09-2024','MM-YYYY'), TO_DATE('11-2026','MM-YYYY'));
