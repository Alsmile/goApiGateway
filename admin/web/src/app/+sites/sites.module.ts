import { NgModule } from '@angular/core';

import { routing } from "./sites.routing";
import { SharedModule } from "../../shared/shared.module";
import {SitesComponent} from './sites.component';
import {SitesHomeComponent} from './home/home.component';
import {SitesHomeService} from './home/home.service';

@NgModule({
  imports:      [ SharedModule, routing ],
  declarations: [
    SitesComponent,
    SitesHomeComponent,
  ],
  providers: [
    SitesHomeService,
  ]
})
export class SitesModule { }
