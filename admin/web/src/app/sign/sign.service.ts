import {Injectable} from '@angular/core';
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
export class SignService {
  public constructor(protected http: HttpService, protected store: StoreService) {
  }

  async Signin(user: any): Promise<any> {
    let ret = await this.http.Post('/api/login', {
      rememberMe: user.rememberMe,
      profile: {email: user.profile.email},
      password: user.password,
      captcha: user.captcha
    });
    if (ret.error) return {};
    if (user.rememberMe) {
      localStorage.setItem("rememberMe", "1");
      localStorage.setItem("token", ret.token);
    }
    else {
      localStorage.removeItem("rememberMe");
      CookieService.set("token", ret.token);
    }
    localStorage.setItem("last.email", user.profile.email);
    this.store.set('user', ret);

    return ret;
  }

  async Signup(user: any): Promise<any> {
    let ret = await this.http.Post('/api/signup', {
      rememberMe: user.rememberMe,
      profile: {email: user.profile.email},
      password: user.password,
      captcha: user.captcha
    });
    if (!ret || ret.error) return false;
    return true;
  }

  async ForgetPassword(user: any): Promise<any> {
    let ret = await this.http.Post('/api/forget/password', user);
    if (!ret || ret.error) return false;
    return true;
  }

  async SignActive(activeCode: string): Promise<any> {
    let ret = await this.http.Post('/api/sign/active', {active:{code: activeCode}});
    if (!ret || ret.error) return false;

    localStorage.setItem("token", ret.token);
    this.store.set('user', ret);
    return true;
  }

  async NewPassword(params: any): Promise<any> {
    let ret = await this.http.Post('/api/sign/new/password', params);
    if (!ret || ret.error) return false;
    return true;
  }

  async SendActiveEmail(email: string): Promise<any> {
    let ret = await this.http.Post('/api/sign/send/active/email', {email: email});
    if (!ret || ret.error) return false;
    return true;
  }
}
