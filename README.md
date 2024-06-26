# CusCus #

Simple full-stack twitter clone (kinda) using HTMX, Go and Postgres.

## Security notes ##

DO NOT USE THIS PROJECT FOR COMMERCIAL PURPOSES AS IS!!! The security is not the best.
While I tried to make it secure from SQL injections (using SQL parameters), XSS (escaping the text on the server),
rooting (the passwords are salted and hashed on the database using a psql extension called crypt), and it uses
https to hopefully avoid MITM attacks, it probably has other vulnerabilities.
I am not a cybersecurity expert and this is just a fun project.

DON'T use your real password on this site. You never know.

## How to run ##

> [!WARNING]
> This has only been tested on Linux (Ubuntu and Arch). Support on other platforms is not guaranteed.

> [!TIP]
> If you are using Mac or Windows, you can install a package manager
> (like chocolatey for Windows or homebrew for Mac), that will make your life easier
> (Or use WSL if on Windows).

### Step 1: Set up Postgres and Go ###

#### Installation ####

On Ubuntu:
```
sudo apt install postgresql golang openssl
```

On Arch
```
sudo pacman -S postgresql go openssl
```
or
```
yay -S postgresql go openssl
```

#### Start Database ####

> [!WARNING]
> This might not work for you. If it doesn't, just paste the error into google. fixing it
> shouldn't be difficult. Remember this is Linux-only,
> for other machines you can just look up how to start the server online.

```
sudo -u postgres initdb --locale=C.UTF-8 --encoding=UTF8 -D '/var/lib/postgres/data'
sudo systemctl enable --now postgresql.service
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

Then connect to the database using:
```
\c cuscus
```

Then finally:

```
CREATE EXTENSION pgcrypto; -- Add pgcrypto extension (for hashing and salting passwords)
CREATE TABLE users (username VARCHAR(40) PRIMARY KEY, password_hash VARCHAR(100), userid SERIAL, sessionid CHAR(64)); -- Creates user table
CREATE TABLE messages (message VARCHAR(10000), username VARCHAR(40), messageid SERIAL PRIMARY KEY); -- Creates messages table (username refers to the sender's username)
```

Now you can exit psql.

```
\q
```

### Step 2: Configure db.go ###

You can skip this step if your password is "fottutapassword".

Open up the "db.go" file in your favourite text
editor and modify the line where it says ```const DB_PWD = "fottutapassword" ```
after the imports, and modify the value to match your password.

### Step 3: Set up TLS certificates ###

Before we can use https we need a certificate.
First, cd into the "certs" directory
```
mkdir certs
cd certs
```
Then, run the following command to generate a private key file and a certificate signing request.
```
openssl req  -new  -newkey rsa:2048  -nodes  -keyout localhost.key  -out localhost.csr
```
It's going to ask you a few questions. You don't need to answer them.

Now that we have those, we need to run the following command to generate the certificate based on them.
```
openssl  x509  -req  -days 365  -in localhost.csr  -signkey localhost.key  -out localhost.crt
```

Now you have 3 files: localhost.csr, localhost.crt and localhost.key. Those are going to be used by go to encrypt the data.

### Step 4: Run the server ###

Go back to the main dir
```
cd ..
```

Type (in the main dir):
```
make bar
```

This will build and run your application. Build files are in the "build" directory.

> [!WARNING]
> Only run the app from the main (cuscus) directory.
> Do not run from the build dir or any other dir, as the file paths in the app are calculated
relative to the current directory, not the binary directory or an absolute path (should be fixed).

## All done!! ##

Now you can visit https://localhost:8080 to login the website and start typing!
