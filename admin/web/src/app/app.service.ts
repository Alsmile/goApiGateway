import {Injectable} from '@angular/core';
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {CookieService, StoreService} from 'le5le-store';

import { HttpService } from '../core/http.service';

@Injectable()
export class AppService extends HttpService {
  public constructor(protected http: Http, protected store: StoreService) {
    super(http, store);
  }

  GetSignConfig(): Observable<any> {
    return this.Get('/api/sign/config');
  }

}
