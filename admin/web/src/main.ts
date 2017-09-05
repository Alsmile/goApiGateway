import { platformBrowserDynamic } from '@angular/platform-browser-dynamic';

import { AppModule } from './app/app.module';
import { enableProdMode } from '@angular/core';


if (process.env.NODE_ENV) enableProdMode();

if (process.env.NODE_ENV === 'production') {
  (<any>window).token = 'token';
} else if (process.env.NODE_ENV === 'dev') {
  (<any>window).token = 'token_dev';
} else {
  (<any>window).token = 'token_dev';
}

platformBrowserDynamic().bootstrapModule(AppModule);

require('./assets/css/app.pcss');
