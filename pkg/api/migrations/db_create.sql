CREATE DATABASE sandpiper;
CREATE USER sandpiper WITH ENCRYPTED PASSWORD 'autocare';
GRANT ALL PRIVILEGES ON DATABASE sandpiper TO sandpiper;