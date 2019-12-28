/*
 * Database: sandpiper
 * Migration: Fill the initial tables with minimal required data
 * Direction: up
 * Date: 2019-11-16
 */
 
BEGIN;

INSERT INTO companies (id, name, active, created_at, updated_at, deleted_at)
    VALUES ('9ad07234-2742-40eb-9ef1-800c2f1164ce', 'sandpiper owner', true,  now(), now(), NULL);

INSERT INTO users (
    id, first_name, last_name, username, password, email, phone, active, role, company_id, created_at, updated_at)
VALUES (
    1, 'Admin', 'Admin', 'admin', '$2a$10$jlCjsODRzvR3L4NVAI6z3uULAAaqMOSadVRWaUJjwtp.cK1huLFlK', 'johndoe@mail.com',
    '314-555-1212 x111', true, 100, '9ad07234-2742-40eb-9ef1-800c2f1164ce', now(), now());

COMMIT;