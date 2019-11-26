/*
 * Database: sandpiper
 * Migration: Create initial database tables in an empty database
 * Direction: Down
 * Date: 2019-11-16
 */
 
BEGIN; 

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS companies;

COMMIT;