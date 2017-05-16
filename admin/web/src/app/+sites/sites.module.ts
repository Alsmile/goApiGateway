import { NgModule } from '@angular/core';

import { routing } from "./sites.routing";
import { SharedModule } from "../../shared/shared.module";
import {SitesComponent} from './sites.component';
import {SitesHomeComponent} from './home/home.component';
import {SitesEditComponent} from './edit/edit.component';
import {SitesApisListComponent} from './apis/list.component';

import {SitesService} from './sites.service';

@NgModule({
  imports:      [ SharedModule, routing ],
  declarations: [
    SitesComponent,
    SitesHomeComponent,
    SitesEditComponent,
    SitesApisListComponent
  ],
  providers: [
    SitesService
  ]
})
export class SitesModule { }
