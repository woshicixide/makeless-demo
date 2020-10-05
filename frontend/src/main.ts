import Makeless from '@makeless/makeless-ui/src/makeless';

// packages
import Config from "@makeless/makeless-ui/src/packages/config/basic/config";
import Router from "@makeless/makeless-ui/src/packages/router/basic/router";
import Axios from "@makeless/makeless-ui/src/packages/http/axios/axios";
import I18n from "@makeless/makeless-ui/src/packages/i18n/basic/i18n";
import LocalStorage from "@makeless/makeless-ui/src/packages/storage/local-storage/local-storage";
import Event from "@makeless/makeless-ui/src/packages/event/basic/event";
import Security from "@makeless/makeless-ui/src/packages/security/basic/security";

// scss
import './scss/app.scss'

// config
import configuration from './../makeless.json';

const config = new Config(configuration);
const router = new Router();
const http = new Axios({baseURL: config.getConfiguration().getHost()});
const i18n = new I18n(config.getConfiguration().getLocale());
const storage = new LocalStorage();
const event = new Event(config.getConfiguration().getHost());
const security = new Security(router, http, event, storage);

new Makeless(config, router, http, i18n, event, security)
    .init()
    .then(makeless => makeless.run());