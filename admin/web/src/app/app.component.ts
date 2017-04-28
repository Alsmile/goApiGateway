import {Component, OnInit} from '@angular/core';
import {Router, ActivatedRoute} from '@angular/router';
import {StoreService} from 'le5le-store';

import {NoticeService} from 'le5le-components';

import {SignService, SignType} from './sign/sign.service'

@Component({
  selector: 'app',
  templateUrl: "app.component.html"
})
export class AppComponent implements OnInit {
  user: any;
  options: any = {showSign: null};
  email: string;
  newPasswordCode: string;
  myStyle: any = {};
  constructor(private _router: Router, private _activateRoute: ActivatedRoute, private _storeService: StoreService,
              private _signService: SignService) {
    this.myStyle = {
      'margin-top': '.6rem',
      'min-height': (document.documentElement.clientHeight - 120) + 'px'
    };
  }

  ngOnInit() {
    // 监听用户信息
    this._storeService.get$('user').subscribe(
      ret => {
        this.user = ret;
      }
    );

    // 检查路由参数判断是否激活账户、找回密码等
    this._activateRoute.queryParams.subscribe((params: any) => {
      if (params) {
        // 激活账户
        if (params.active) {
          this._signService.SignActive(params.active).subscribe(
            ret => {
              let _noticeService: NoticeService = new NoticeService();
              _noticeService.notice({body: '恭喜账号激活成功，感谢您的使用！', theme: 'success'});
            },
            err => {
              let _noticeService: NoticeService = new NoticeService();
              _noticeService.notice({body: '激活失败，请稍后重试！', theme: 'error', timeout: 5000});
            }
          );
        }
        // 找回密码方式设置新密码
        else if (params.forgetPassword) {
          this.newPasswordCode = params.forgetPassword;
          this.email = params.email;
          this.options.showSign = SignType.NewPasswordDialog;
        }
      }
    });

    // 需要激活
    this._storeService.get$('needActive').subscribe(
      ret => {
        if (ret) this.options.showSign = SignType.NeedActiveDialog;
      }
    );
  }

  isActive (strUrl: string) {
    if (!strUrl || strUrl === "/") {
      return !this._router.url || this._router.url === '/';
    }
    else {
      return this._router.url.indexOf(strUrl) === 0;
    }
  }

  onSignin () {
    this.options.showSign = SignType.SignInDialog;
  }

  onSignup () {
    this.options.showSign = SignType.SignUpDialog;
  }

  onSignout () {
    this.options.showSign = null;
    this._storeService.set('auth', null);
  }

  onSignCancel (event:any) {
    if (event) event.stopPropagation();

    if (this.options.showSign === SignType.NewPasswordDialog) this._router.navigate(['/']);
    this.options.showSign = null;
  }

  onSendActiveEmail (event:any) {
    if (event) event.stopPropagation();

    if (!this.user) return;

    this._signService.SendActiveEmail(this.user.email).subscribe(
      ret => {
        let _noticeService: NoticeService = new NoticeService();
        _noticeService.notice({body: '已经发送激活邮件到您的邮箱，请注意查收！', theme: 'success'});
      },
      err => {
        let _noticeService: NoticeService = new NoticeService();
        _noticeService.notice({body: '发送激活邮件失败，请稍后重试！', theme: 'error', timeout: 5000});
      }
    );
  }
}
