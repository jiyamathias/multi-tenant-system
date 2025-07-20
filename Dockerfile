FROM golang:alpine
# RUN apk add git

ARG app_env
ENV APP_ENV $app_env

WORKDIR /go/src/codematic/src

COPY . /go/src/codematic
RUN touch /tmp/runner-build-errors.log 
# Set the environment variable to use the .env file
ENV DOTENV_PATH=../../.env
#RUN go mod download
RUN go get ./
# Remove .git directory
RUN rm -rf .git
RUN go build
RUN go install github.com/pilu/fresh@latest

# if dev setting will use pilu/fresh for code reloading via docker-compose volume sharing with local machine
# if production setting will build binary
CMD if [ ${APP_ENV} = production ]; \
	then \
	src; \
	else \
	fresh; \
	fi

EXPOSE 5002