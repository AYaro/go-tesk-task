BEGIN;

CREATE TABLE part_manufacturer (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  created_at  TIMESTAMP default CURRENT_TIMESTAMP
);

CREATE TABLE part(
  id serial PRIMARY KEY,
  manufacturer_id int REFERENCES part_manufacturer(id) ON DELETE CASCADE,
  vendor_code VARCHAR(100),
  created_at  TIMESTAMP default CURRENT_TIMESTAMP,
  deleted_at  TIMESTAMP
);

INSERT INTO part_manufacturer (name)
VALUES ('Manufacturer 1');
INSERT INTO part_manufacturer (name)
VALUES ('Manufacturer 2');
INSERT INTO part_manufacturer (name)
VALUES ('Tesla');
INSERT INTO part_manufacturer (name)
VALUES ('Volkswagen');
INSERT INTO part_manufacturer (name)
VALUES ('Toyota');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (1, 'some vendor code');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (1, 'some vendor code 2');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (2, 'test vendor code');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (3, '1526GHKL13');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (3, 'HJKH167649');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (3, '978973HJSD');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (4, '123');
INSERT INTO part (manufacturer_id, vendor_code)
VALUES (5, 'test');

COMMIT;