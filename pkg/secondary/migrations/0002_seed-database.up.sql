/*
 * Database: sandpiper
 * Migration: Fill the initial tables with minimal required data
 * Direction: up
 * Date: 2019-11-16
 */
 
BEGIN;

INSERT INTO companies (id, name, active, created_at, updated_at)
    VALUES (1, 'admin_company', true,  now(), now());

INSERT INTO users (
    id, first_name, last_name, username, password, email, active, role, company_id, created_at, updated_at)
VALUES (
    1, 'Admin', 'Admin', 'admin', '$2a$10$jlCjsODRzvR3L4NVAI6z3uULAAaqMOSadVRWaUJjwtp.cK1huLFlK', 'johndoe@mail.com',
    true, 100, 1, now(), now());

COMMIT;