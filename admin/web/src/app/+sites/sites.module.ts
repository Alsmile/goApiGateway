import { NgModule } from '@angular/core';

import { routing } from "./sites.routing";
import { SharedModule } from "../../shared/shared.module";
import {SitesComponent} from './sites.component';
import {SitesHomeComponent} from './home/home.component';
import {SitesHomeService} from './home/home.service';
import {SitesEditComponent} from './edit/edit.component';
import {SitesEditService} from './edit/edit.service';

@NgModule({
  imports:      [ SharedModule, routing ],
  declarations: [
    SitesComponent,
    SitesHomeComponent,
    SitesEditComponent
  ],
  providers: [
    SitesHomeService,
    SitesEditService
  ]
})
export class SitesModule { }
