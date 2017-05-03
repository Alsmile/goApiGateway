import {EventEmitter, Injectable} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Http} from "@angular/http";
import {CookieService, StoreService} from 'le5le-store';

import 'rxjs/add/operator/do';

import { HttpService } from '../../core/http.service';

export enum SignType {
  SignInDialog = 1,
  SignUpDialog,
  ForgetPasswordDialog,
  NeedActiveDialog,
  NewPasswordDialog,
}

const Salt:string = 'goMicroServer.io';

@Injectable()
export class SignService extends HttpService {
  public constructor(protected http: Http, protected store: StoreService) {
    super(http, store);
  }

  Signin(user: any): Observable<any> {
    let u: any = {
      rememberMe: user.rememberMe,
      email: user.email,
      password: user.password,
      captcha: user.captcha
    };
    let login$ = this.Post('/api/login', u).do( ret => {
      if (ret.error) return;
      if (user.rememberMe) {
        localStorage.setItem("rememberMe", "1");
        localStorage.setItem("token", ret.token);
      }
      else {
        localStorage.removeItem("rememberMe");
        CookieService.set("token", ret.token);
      }
      localStorage.setItem("last.email", user.email);
      this.store.set('user', ret);
    });
    return login$;
  }

  Signup(user: any): Observable<any> {
    let u: any = {
      rememberMe: user.rememberMe,
      email: user.email,
      password: user.password,
      captcha: user.captcha
    };

    return this.Post('/api/signup', u);
  }

  ForgetPassword(user: any): Observable<any> {
    return this.Post('/api/forget/password', user);
  }

  SignActive(activeCode: string): Observable<any> {
    let login$ = this.Post('/api/sign/active', {activeCode: activeCode}).do( ret => {
      localStorage.setItem("token", ret.token);
      this.store.set('user', ret);
    });
    return login$;
  }

  NewPassword(params: any): Observable<any> {
    return this.Post('/api/sign/new/password', params);
  }

  SendActiveEmail(email: string): Observable<any> {
    return this.Post('/api/sign/send/active/email', {email: email});
  }
}
