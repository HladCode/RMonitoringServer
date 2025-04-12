--CREATE DATABASE sensor_data;
--\c sensor_data

-- CREATE EXTENSION IF NOT EXISTS timescaledb;

-- CREATE TABLE IF NOT EXISTS sensor_readings (
--     timestamp         TIMESTAMPTZ       NOT NULL,
--     device_id         TEXT              NOT NULL,
--     sensor_index      INTEGER           NOT NULL,
--     -- measurement_type  TEXT              NOT NULL,
--     value             DOUBLE PRECISION  NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS sensor_index_meaning (
--     device_id         TEXT              NOT NULL,
--     sensor_index      INTEGER           NOT NULL,
--     measurement_type  TEXT              NOT NULL,
--     meaning           TEXT              NOT NULL  -- щоб розуміти що саме вимірює той чи інший датчик з певним індексом
-- );

-- SELECT create_hypertable('sensor_readings', 'timestamp');

-- CREATE INDEX idx_device_id ON sensor_readings (device_id);
-- CREATE INDEX idx_sensor_index ON sensor_readings (sensor_index);
-- --CREATE INDEX idx_measurement_type ON sensor_readings (measurement_type);
-- 
--
-- Все що закоментовано має бути виконано на початку створення бази даних.
--
-- TODO: в методі бд Init створити перевірку на наявність таблиць, якщо вони відстутні,
-- що буде означати, що запитів зверху не було, то виконати код зверху!! 
--

CREATE OR REPLACE FUNCTION get_sensor_data_for_day(
    p_device_id TEXT,
    p_sensor_index INT,
    p_date DATE
)
RETURNS TABLE (
    "timestamp" TIMESTAMPTZ,
    value DOUBLE PRECISION
) AS $$
BEGIN
    RETURN QUERY
    SELECT sensor_readings."timestamp", sensor_readings.value
    FROM sensor_readings
    WHERE device_id = p_device_id
      AND sensor_index = p_sensor_index
      AND sensor_readings."timestamp" >= p_date
      AND sensor_readings."timestamp" < p_date + INTERVAL '1 day'
    ORDER BY sensor_readings."timestamp";
END;
$$ LANGUAGE plpgsql;