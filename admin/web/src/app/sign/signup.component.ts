import { Component, Input, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';

import { SignService, SignType } from "./sign.service";

@Component({
  selector: 'sign-up',
  templateUrl: "signup.component.html"
})
export class SignupComponent {
  captchaUrl: string;
  user: any = {profile:{email: ''}, password: '', captcha: ''};
  @Input() options: any;
  saving: boolean;
  success: boolean;
  @ViewChild('signupForm') currentForm: NgForm;
  constructor(private _signService: SignService) {
  }

  ngOnInit() {
    this.getCaptcha();
  }

  getCaptcha() {
    this.captchaUrl = "/captcha?rand="+ new Date().getTime();
  }

  onSignin() {
    if (this.options) this.options.showSign = SignType.SignInDialog;
  }

  onCancel() {
    if (this.options) this.options.showSign = null;
  }

  onSubmit(valid: boolean) {
    if (this.currentForm.form.invalid) return;

    this.saving = true;
    this._signService.Signup(this.user).subscribe(
      ret => {
        this.success = true;
      },
      err => console.error(err),
      () => this.saving = false
    );
  }
}
