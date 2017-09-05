import { ModuleWithProviders, NgModule, Optional, SkipSelf } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';

import { Le5leStoreModule, StoreService, CookieService } from 'le5le-store';
import { NoticeService } from 'le5le-components';
import { HttpService } from './http.service';

@NgModule({
  imports:       [ CommonModule, Le5leStoreModule ],
  declarations: [ ],
  exports:       [ ],
  providers:     [
    NoticeService,
    HttpService,
  ]
})
export class CoreModule {
  constructor(@Optional() @SkipSelf() parentModule: CoreModule,
              private _router: Router, private _storeService: StoreService, private _httpService: HttpService) {

    if (parentModule) {
      throw new Error(
        'CoreModule is already loaded. Import it in the AppModule only');
    }

    this._storeService.set('author', 'alsmile');

    // 监听用户认证
    this._storeService.get$('auth').subscribe(
      ret => {
        // 认证失败
        if (!ret || ret === -1) {
          this._storeService.set('user', null);
          localStorage.removeItem((<any> window).token);
          CookieService.delete((<any> window).token, {domain: document.domain.split('.').slice(-2).join('.')});

          let loginUrl = this._storeService.get('loginUrl');
          if (loginUrl) return window.location.href = loginUrl + window.location.href;

          this._router.navigate(['/']);

          if (!ret) {
            this._storeService.set('needActive', 1);
          }
        }
      }
    );

    this.onProfile();
  }

  async onProfile(): Promise<void> {
    let ret = await this._httpService.Get('/api/user/profile');
    if (ret && !ret.error) this._storeService.set('user', ret);
    else this._storeService.set('auth', -1);
  }

  static forRoot(): ModuleWithProviders {
    return {
      ngModule: CoreModule,
      providers: [ ]
    };
  }
}
