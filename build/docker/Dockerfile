ARG GO_IMAGE

FROM ${GO_IMAGE} AS base

ARG PORT
ARG APP_NAME

ENV APP_NAME=${APP_NAME}
ENV PORT=${PORT}

WORKDIR /${APP_NAME}
COPY . .

FROM base AS local

WORKDIR /
RUN go install github.com/cosmtrek/air@latest

WORKDIR /${APP_NAME}

CMD air -c ./configs/.air.toml

EXPOSE ${PORT}

