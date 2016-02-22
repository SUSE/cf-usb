# mysql

https://hub.docker.com/r/mysql/mysql-server/

To start the container:
```
docker run -p 3306:3306 --name mysql_container -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql/mysql-server:5.6
```
Most of the variables listed below are optional, but one of the variables MYSQL_ROOT_PASSWORD, MYSQL_ALLOW_EMPTY_PASSWORD, MYSQL_RANDOM_ROOT_PASSWORD must be given.

MYSQL_ROOT_PASSWORD
This variable specifies a password that will be set for the MySQL root superuser account. NOTE: Setting the MySQL root user password on the command line is insecure. 

MYSQL_RANDOM_ROOT_PASSWORD
When this variable is set to yes, a random password for the server's root user will be generated. The password will be printed to stdout in the container, and it can be obtained by using the command docker logs my-container-name.

MYSQL_ONETIME_PASSWORD
This variable is optional. When set to yes, the root user's password will be set as expired, and must be changed before MySQL can be used normally. This is only supported by MySQL 5.6 or newer.

MYSQL_DATABASE
This variable is optional. It allows you to specify the name of a database to be created on image startup. If a user/password was supplied (see below) then that user will be granted superuser access (corresponding to GRANT ALL) to this database.

MYSQL_USER, MYSQL_PASSWORD
These variables are optional, used in conjunction to create a new user and set that user's password. This user will be granted superuser permissions (see above) for the database specified by the MYSQL_DATABASE variable. Both variables are required for a user to be created.

Do note that there is no need to use this mechanism to create the root superuser, that user gets created by default with the password set by either of the mechanisms (given or generated) discussed above.

MYSQL_ALLOW_EMPTY_PASSWORD
Set to yes to allow the container to be started with a blank password for the root user. NOTE: Setting this variable to yes is not recommended unless you really know what you are doing, since this will leave your MySQL instance completely unprotected, allowing anyone to gain complete superuser access.



# postgresql

https://hub.docker.com/_/postgres/

To start the container :
```
docker run -p 5432:5432 --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres:9.4
```
The PostgreSQL image uses several environment variables which are easy to miss. While none of the variables are required, they may significantly aid you in using the image.

POSTGRES_PASSWORD
This environment variable is recommended for you to use the PostgreSQL image. This environment variable sets the superuser password for PostgreSQL. The default superuser is defined by the POSTGRES_USER environment variable. In the above example, it is being set to "mysecretpassword".

POSTGRES_USER
This optional environment variable is used in conjunction with POSTGRES_PASSWORD to set a user and its password. This variable will create the specified user with superuser power and a database with the same name. If it is not specified, then the default user of postgres will be used.

PGDATA
This optional environment variable can be used to define another location - like a subdirectory - for the database files. The default is /var/lib/postgresql/data, but if the data volume you're using is a fs mountpoint (like with GCE persistent disks), Postgres initdb recommends a subdirectory (for example /var/lib/postgresql/data/pgdata ) be created to contain the data.

POSTGRES_DB
This optional environment variable can be used to define a different name for the default database that is created when the image is first started. If it is not specified, then the value of POSTGRES_USER will be used.



# mongo

https://hub.docker.com/_/mongo/

There are several ways to store data used by applications that run in Docker containers.
 
Let Docker manage the storage of your database data by writing the database files to disk on the host system using its own internal volume management. This is the default and is easy and fairly transparent to the user. The downside is that the files may be hard to locate for tools and applications that run directly on the host system, i.e. outside containers.

Create a data directory on the host system (outside the container) and mount this to a directory visible from inside the container. This places the database files in a known location on the host system, and makes it easy for tools and applications on the host system to access the files. The downside is that the user needs to make sure that the directory exists, and that e.g. directory permissions and other security mechanisms on the host system are set up correctly.

WARNING: because MongoDB uses memory mapped files it is not possible to use it through vboxsf to your host (vbox bug). VirtualBox shared folders are not supported by MongoDB (see docs.mongodb.org and related jira.mongodb.org bug). This means that it is not possible with the default setup using Docker Toolbox to run a MongoDB container with the data directory mapped to the host.


To start the container :
```
docker run -p 27017:27017 --name some-mongo -d mongo:3.1
```
To start the container with a shared folder :
```
docker run -p 27017:27017 --name some-mongo -v /my/own/datadir:/data/db -d mongo:3.1
```

# redis
https://hub.docker.com/_/redis/

You can create your own Dockerfile that adds a redis.conf from the context into /data/, like so.
```
FROM redis
COPY redis.conf /usr/local/etc/redis/redis.conf
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
```
or:
```
docker run -v /myredis/conf/redis.conf:/usr/local/etc/redis/redis.conf --name myredis redis redis-server /usr/local/etc/redis/redis.conf
```

To start the container :
```
docker run -p 6379:6379 --name some-redis -d redis:3
```

