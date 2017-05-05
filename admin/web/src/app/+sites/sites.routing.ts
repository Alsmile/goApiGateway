import { ModuleWithProviders }  from '@angular/core';
import { RouterModule } from '@angular/router';

import {SitesComponent} from "./sites.component";
import {SitesHomeComponent} from "./home/home.component";
import {SitesEditComponent} from './edit/edit.component';


export const routing: ModuleWithProviders = RouterModule.forChild([
  {
    path: '',
    component: SitesComponent ,
    children: [
      { path: 'home', component: SitesHomeComponent },
      { path: 'edit', component: SitesEditComponent },
    ]
  },
]);
