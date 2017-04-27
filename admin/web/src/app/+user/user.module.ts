import { NgModule } from '@angular/core';

import { userRouting } from "./user.routing";
import { SharedModule } from "../../shared/shared.module";
import { UserComponent } from "./user.component";
import { UserHomeComponent } from "./home/home.component"
import { UserHomeService } from "./home/home.service"

@NgModule({
  imports:      [ SharedModule, userRouting ],
  declarations: [
    UserComponent,
    UserHomeComponent,
  ],
  providers: [
    UserHomeService
  ]
})
export class UserModule { }
