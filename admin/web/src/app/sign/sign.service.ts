import {EventEmitter, Injectable} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Http} from "@angular/http";
import {CookieService, StoreService} from 'le5le-store';
import HmacSHA256 from 'crypto-js/hmac-sha256';

import 'rxjs/add/operator/do';

import { HttpService } from '../../core/http.service';

export enum SignType {
  SignInDialog = 1,
  SignUpDialog,
  ForgetPasswordDialog,
  NeedActiveDialog,
  NewPasswordDialog,
}

const Salt:string = 'le5le.com';

@Injectable()
export class SignService extends HttpService {
  public constructor(protected http: Http, protected store: StoreService) {
    super(http, store);
  }

  Signin(user: any): Observable<any> {
    let pwd:string = HmacSHA256(user.password, Salt)+'';
    pwd = HmacSHA256(pwd, user.email)+'';
    let u: any = {
      rememberMe: user.rememberMe,
      email: user.email,
      password: pwd,
      captcha: user.captcha
    };
    let login$ = this.Post('/api/login', u).do( ret => {
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
    let pwd:string = HmacSHA256(user.password, Salt)+'';
    let u: any = {
      rememberMe: user.rememberMe,
      email: user.email,
      password: pwd,
      captcha: user.captcha
    };

    return this.Post('/api/signup', u);
  }

  ForgetPassword(user: any): Observable<any> {
    return this.Post('/api/forget/password', user);
  }

  SignActive(activeCode: string): Observable<any> {
    let login$ = this.Post('/api/sign/active', {code: activeCode}).do( ret => {
      localStorage.setItem("token", ret.token);
      this.store.set('user', ret);
    });
    return login$;
  }

  NewPassword(params: any): Observable<any> {
    params.password = HmacSHA256(params.password, Salt)+'';
    return this.Post('/api/sign/new/password', params);
  }

  SendActiveEmail(email: string): Observable<any> {
    return this.Post('/api/sign/send/active/email', {email: email});
  }
}
