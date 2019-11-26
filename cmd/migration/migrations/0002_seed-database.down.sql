/*
 * Database: sandpiper
 * Migration: Fill the initial tables with minimal required data
 * Direction: down
 * Date: 2019-11-16
 */
 
BEGIN;

DELETE FROM users;
DELETE FROM roles;
DELETE FROM locations;
DELETE FROM companies;

COMMIT;