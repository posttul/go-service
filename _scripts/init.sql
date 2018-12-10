CREATE TABLE lot (
    id SERIAL PRIMARY KEY,
    name text,
    address text,
    country text
);
INSERT INTO "lot"("id","name","address","country") VALUES
(1,E'Peten14',E'Peten 14  Col. Narvarte Poniente Benito Juarez CDMX', E'mx');
