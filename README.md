# Makeless Demo

SaaS Framework - Production-Ready Docker Demo - based on [Google Distroless](https://github.com/GoogleContainerTools/distroless) and [Apache](https://hub.docker.com/_/httpd)

[![Build Status](https://ci.loeffel.io/api/badges/makeless/makeless-go/status.svg)](https://ci.loeffel.io/makeless/makeless-demo) 

- Frontend: [makeless-ui](https://github.com/makeless/makeless-ui)
- Backend: [makeless-go](https://github.com/makeless/makeless-go)

## Run: localhost

- Configure your mailer in your [docker-compose.yml](https://github.com/makeless/makeless-demo/blob/master/docker-compose.yml#L29) file 
- Run `docker-compose up -d --build`
