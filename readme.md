
# JTI Backend Test
API Docs: you can import api-docs.json to postman.

You can run this code by 2 way:

## Run with docker-compose

#### 1. Adjust .env
Make sure to adjust file .env before you build docker image or you can edit from docker-compose.yml (services.jti.environment).

#### 2. Build docker image
```
docker build -t jtitest:latest .
```

#### 3. Run docker-compose
```
docker compose up -d
```

#### 4. Import database

You can import database jti_test.sql from phpmyadmin in http://localhost:8081. Reference: https://help.one.com/hc/en-us/articles/115005588189-How-do-I-import-a-database-to-phpMyAdmin.

#### 5. Open website

Open the website (http://localhost:8082) on your favorite browser and make sure you allow sound permission to got sound notification. Reference: https://help.taskworld.com/en/articles/4841740-how-to-allow-or-block-sites-to-play-sound-in-google-chrome. 

## Run manually

#### 1. Import database
You can find the file on this repo named jti_test.sql.

#### 2. Adjust .env
Make sure to adjust file .env before you run the code.

#### 3. Run the code
Please ensure that the environment variables in .env are correct.
```
go run main.go
```
#### 4. Open website

Open the website (http://localhost:8082) on your favorite browser and make sure you allow sound permission to got sound notification. Reference: https://help.taskworld.com/en/articles/4841740-how-to-allow-or-block-sites-to-play-sound-in-google-chrome.





