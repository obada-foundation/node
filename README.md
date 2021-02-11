# Contains Five Dockerized Containers.

* gateway -- Nginx and Vue.js front-end, with PHP Laravel / API business logic to connect to the client via the OBADA API 
* db -- a single MySql database table storing a synchronized “view” of the blockchain.
* QLDB  -- QLDB-specific blockchain drivers, written in Go Language.
* queue -- A messaging queue (Amazon SQS) to synchrnoize the obits table with QLDB.
* zipkin -- a distributed tracing utility for debugging and error correction.

