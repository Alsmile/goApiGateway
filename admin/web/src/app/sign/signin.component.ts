import { Component, Input, ViewChild } from '@angular/core';
import {Router} from '@angular/router';
import { NgForm } from '@angular/forms';

import { SignService, SignType } from "./sign.service";


@Component({
  selector: 'sign-in',
  templateUrl: "signin.component.html"
})
export class SigninComponent {
  captchaUrl: string = '';
  user:any = {profile:{email: localStorage.getItem("last.email") || ''}, password: ''};
  @Input() options: any;
  saving: boolean;
  @ViewChild('signinForm') currentForm: NgForm;
  constructor(private _signService: SignService,  private _router: Router) {
  }

  getCaptcha() {
    this.captchaUrl = "/captcha?rand="+ new Date().getTime();
  }

  onSignup() {
    if (this.options) this.options.showSign = SignType.SignUpDialog;
  }

  onForgetPassword() {
    if (this.options) this.options.showSign = SignType.ForgetPasswordDialog;
  }

  async onSubmit(): Promise<void> {
    if (this.currentForm.form.invalid) return;

    this.saving = true;
    let ret = await this._signService.Signin(this.user);
    if (ret.errorTip) {
      this.getCaptcha();
    } else {
      this._router.navigateByUrl('/sites/home');
    }
    this.saving = false;
  }
}
