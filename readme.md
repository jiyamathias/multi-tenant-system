## Multi-Tenant API

##### Things to know and to get started with it's engineering.

- This project runs with a docker. You only need to execute `docker-compose up` from your favorite terminal. Keep in mind that the first time, it will run for a while, and download some stuff, stay cool!
- Golang is used with its the most amazing gin framework for http handling. APIs are in REST.

##### So, welcome to Go!.

The main source code is in the `src` directory. Don't be scared!<br/>
###### Once docker is running, you can always access the app via the follow:<br/>

- Health check: `http://localhost:5002`
- Rest API endpoint: `http://localhost:5002/api/v1`
- API docs endpoint: `http://localhost:5002/api/v1/docs/index.html`

Backend is connected and using `postgreSQL` and `Redis` which is already up once docker compose is running.
Credentials can be found in `src/.env` file.
In development, you would need to create a dev.env file because this would be where the `env` configuration values would be gotten from.
NOTE: running `docker compose up` would return an error until the dev.env file is created.

For running tests, you can run the following command `go test ./...`

We use GORM's(`https://gorm.io`) `AutoMigrate()` to automatically make migrations. You can check the `./src/storage/storage.go` 

### Project breakdown
You'll find in the `/src` folder the main source code running this application
#### `/src/controller`
The controller package contains files that defines the logic of for certain operations and then makes a database call to the using the appropriate method for that respective operation.

#### `/src/docs`
Do not update this folder manually, it is updated automatically when you run the `swag init` command to generate documentation.

#### `/src/handler`
The handler package contains all the API endpoints exposed from the application
- #### `/src/handler/handler.go` 
This file resgiters all the endpoints created and exposes it.

#### `/src/model`
The model package contains files for respective models

#### `/src/pgk`
The pkg package contains internal methods, utility functions and so on.
- #### `/src/pkg/environment` 
This package contains functions responsible to env variables setup and management.
- #### `/src/pkg/helper` 
This package contains utility(helper) functions.
- #### `/src/pkg/middleware` 
This package contains methods responsible for middlewares.

#### `/src/storage`
The storage package contains files which are named to easily locate  that file that handles a particular database table interaction. What this file does is basically handle database interactions.
- #### `/src/storage/redis` 
This redis package contains methods that handler the redis interaction.

#### `/src/thirdparty`
The thirdparty package contains methods that interact with thirdparty services.
#### `/src/temp`
Do not touch this file. It is created from running the `docker compose up` and logs errors while running the docker file if any.

#### `/src/go.mod` and `/src/go.sum`
Do not touch, this file handles the go modules and dependencies this application depends on.

#### `/src/main.go`
This file is the entrance to the application.