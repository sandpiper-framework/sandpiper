/*
   Copyright The Sandpiper Authors. All rights reserved.
   Use of this source code is governed by an MIT-style
   license that can be found in the LICENSE.md file.
*/

/************************************************************************/
/* primary database (remove this section if not publishing)             */
/************************************************************************/
\set primary_dbname 'sandpiper'
\set primary_dbuser 'sandpiper'
\set primary_dbpass 'autocare'

CREATE DATABASE :primary_dbname;
CREATE USER :primary_dbuser WITH ENCRYPTED PASSWORD :'primary_dbpass';
GRANT ALL PRIVILEGES ON DATABASE :primary_dbname TO :primary_dbuser;

/************************************************************************/
/* secondary database (remove this section if not subscribing)          */
/************************************************************************/
\set secondary_dbname 'tidepool'
\set secondary_dbuser 'sandpiper'
\set secondary_dbpass 'autocare'

CREATE DATABASE :secondary_dbname;
CREATE USER :secondary_dbuser WITH ENCRYPTED PASSWORD :'secondary_dbpass';
GRANT ALL PRIVILEGES ON DATABASE :secondary_dbname TO :secondary_dbuser;
