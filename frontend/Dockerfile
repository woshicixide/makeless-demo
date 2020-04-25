FROM node:13 as build
RUN yarn global add @vue/cli @vue/cli-service-global

WORKDIR /go-saas-ui
ADD . /go-saas-ui

RUN yarn install
RUN yarn build

FROM httpd:2.4
COPY --from=build /go-saas-ui/dist /usr/local/apache2/htdocs/