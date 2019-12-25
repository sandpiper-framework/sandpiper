/*
 * Database: sandpiper
 * Migration: Fill the initial tables with minimal required data
 * Direction: up
 * Date: 2019-11-16
 */
 
BEGIN;

INSERT INTO companies (id, name, active, created_at, updated_at, deleted_at)
    VALUES (1, 'admin_company', true,  now(), now(), NULL);

INSERT INTO locations (id, name, address, active, company_id, created_at, updated_at, deleted_at)
    VALUES (1, 'admin_location', 'admin_address', true, 1, now(), now(), NULL);

INSERT INTO roles (id, access_level, name)
VALUES
    (100, 100, 'SUPER_ADMIN'),
    (110, 110, 'ADMIN'),
    (120, 120, 'COMPANY_ADMIN'),
    (130, 130, 'LOCATION_ADMIN'),
    (200, 200, 'USER');

INSERT INTO users (
    id, first_name, last_name, username, password, email, active, role_id, company_id, location_id, created_at, updated_at)
VALUES (
    1, 'Admin', 'Admin', 'admin', '$2a$10$jlCjsODRzvR3L4NVAI6z3uULAAaqMOSadVRWaUJjwtp.cK1huLFlK', 'johndoe@mail.com',
    true, 100, 1, 1, now(), now());

COMMIT;