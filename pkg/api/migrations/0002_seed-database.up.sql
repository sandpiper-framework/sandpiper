/*
 * Database: sandpiper
 * Migration: Fill the initial tables with minimal required data
 * Direction: up
 * Date: 2019-11-16
 *
 * TODO: Upon start, check for primary company and admin user... prompt for missing information.
 */
 
BEGIN;

INSERT INTO companies
    (id, name, sync_addr, active, created_at, updated_at)
VALUES
    ('10000000-0000-0000-0000-000000000000', 'Primary Company', 'https://sandpiper.primary.com', true,  now(), now());


INSERT INTO users (
    id, first_name, last_name, username, password, email, phone, active, role, company_id, created_at, updated_at)
VALUES
    (DEFAULT, 'Admin', 'Admin', 'admin', '$2a$10$jlCjsODRzvR3L4NVAI6z3uULAAaqMOSadVRWaUJjwtp.cK1huLFlK', 'admin@mail.com',
    '314-555-1212 x111', true, 100, '10000000-0000-0000-0000-000000000000', now(), now());

COMMIT;