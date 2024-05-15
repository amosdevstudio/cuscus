# CusCus #

Simple full-stack project in HTMX, Go and Postgres.

## Security notes ##

DO NOT USE THIS PROJECT FOR COMMERCIAL PURPOSES AS IS!!! The security is very much flawed.
While I tried to make it secure from SQL injections (using SQL parameters), XSS (escaping the text on the server),
and rooting (the passwords are salted and hashed on the database using a psql extention called crypt), it is surely
vulnerable to a simple MITM attack,
as it doen't use https and it sends username and password in PLAIN TEXT to
send any message (as it doesn't use session cookies), and I'm sure there are many more vulnerabilities (I'm not an expert).

## How to run ##

> [!WARNING]
> This has only been tested on Linux (Ubuntu and Arch). Support on other platforms is not guaranteed.

> [!TIP]
> If you are using Mac or Windows, you can install a package manager
(like chocolatey for Windows or homebrew for Mac), that will make your life easier.

### Step 1: Set up Postgres and Go ###

#### Installation ####

On Ubuntu:
```
sudo apt install postgresql golang
```


On Arch
```
sudo pacman -S postgresql go
```
or
```
yay -S postgresql go
```

#### Start Database ####

> [!WARNING]
> This might not work for you. If it doesn't, just paste the error into google. fixing it
> shouldn't be difficult. Remember this is Linux-only,
> for other machines you can just look up how to start the server online.

```
sudo systemctl start postgresql.service
sudo -u postgres psql
```

Yiuppieee!! If everything is right you should see an SQL cli. Otherwise you just got an error :( ...

#### Set postgres password ####

Once you got into the psql cli, you have to set a password for the postgres user.
Remember this password as you will need it later. Type:
```
\password
```

Now you will be prompted to type a new password. You can set your own or you can use mine.
If you use mine you can skip a step later. My password is "fottutapassword", without the quotes :).

#### Create database and tables ####

Just paste this code into your terminal once you got in the psql cli:
```
DROP DATABASE IF EXISTS cuscus; -- Delete cuscus db if exists
CREATE DATABASE cuscus; -- Well... that
```

Then type:
```
\c cuscus
```

Then finally:

```
CREATE EXTENSION pgcrypto; -- Add pgcrypto extension (for hashing and salting passwords)
CREATE TABLE users (username VARCHAR(40) PRIMARY KEY, password_hash VARCHAR(100), userid SERIAL); -- Creates user table
CREATE TABLE messages (message VARCHAR(10000), username VARCHAR(40), messageid SERIAL PRIMARY KEY); -- Creates messages table (username refers to the sender's username)
```

### Step 2: Configure db.go ###

You can skip this step if your password is "fottutapassword".

Open up the "db.go" file in your favourite text
editor and modify the line where it says ```const DB_PWD = "fottutapassword" ```
after the imports, and modify the value to match your password.

### Step 3: Run the server ###

Type (in the main dir):

```
go run *.go
```

## All done!! ##

Now you can visit http://localhost:8080 to login the website and start typing!
