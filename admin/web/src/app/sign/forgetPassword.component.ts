import { Component, Input, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';

import { SignService, SignType } from "./sign.service";

@Component({
  selector: 'forget-password',
  templateUrl: "forgetPassword.component.html"
})
export class ForgetPasswordComponent {
  captchaUrl: string;
  @Input() options: any;
  user: any = {profile:{email: ''}, captcha: ''};
  saving: boolean;
  success: boolean;
  @ViewChild('forgetForm') currentForm: NgForm;
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

  onSubmit () {
    if (this.currentForm.form.invalid) return;

    this.saving = true;
    this._signService.ForgetPassword(this.user).subscribe(
      ret => {
        this.success = true;
      },
      err => console.error(err),
      () => this.saving = false
    );
  }
}
