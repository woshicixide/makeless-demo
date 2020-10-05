FROM node:14 as build

WORKDIR /go-saas-ui
COPY ./makeless.json /go-saas-ui
COPY ./frontend /go-saas-ui

RUN yarn install
RUN yarn build

FROM httpd:2.4
COPY --from=build /go-saas-ui/dist /usr/local/apache2/htdocs/
COPY --from=build /go-saas-ui/.htaccess /usr/local/apache2/htdocs/.htaccess

RUN sed -i '/<Directory "\/usr\/local\/apache2\/htdocs">/,/<\/Directory>/ s/AllowOverride None/AllowOverride all/' conf/httpd.conf

