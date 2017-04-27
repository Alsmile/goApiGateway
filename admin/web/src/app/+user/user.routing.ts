import { ModuleWithProviders }  from '@angular/core';
import { RouterModule } from '@angular/router';

import {UserComponent} from "./user.component";
import {UserHomeComponent} from "./home/home.component";


export const userRouting: ModuleWithProviders = RouterModule.forChild([
  {
    path: '',
    component: UserComponent ,
    children: [
      { path: 'home', component: UserHomeComponent },
    ]
  },
]);
