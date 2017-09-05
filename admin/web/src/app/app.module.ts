import { NgModule }       from '@angular/core';
import { BrowserModule }  from '@angular/platform-browser';
import { HttpModule } from '@angular/http';

import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';
import { AppComponent } from './app.component';
import { routing, appRoutingProviders } from './app.routing';

import { HomeComponent } from './home/home.component';
import { SigninComponent } from './sign/signin.component';
import { SignupComponent } from './sign/signup.component';
import { ForgetPasswordComponent } from './sign/forgetPassword.component';
import { NewPasswordComponent } from './sign/newPassword.component';

import { SignService } from './sign/sign.service';

@NgModule({
  imports: [
    BrowserModule,
    HttpModule,
    CoreModule.forRoot(),
    SharedModule,
    routing
  ],
  declarations: [
    AppComponent,
    HomeComponent,
    SigninComponent,
    SignupComponent,
    ForgetPasswordComponent,
    NewPasswordComponent,
  ],
  providers: [
    SignService,
    appRoutingProviders
  ],
  bootstrap: [ AppComponent ]
})
export class AppModule {
}
