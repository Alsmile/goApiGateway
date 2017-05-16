import {Injectable} from '@angular/core';
import {Observable} from "rxjs/Observable";
import {Http} from "@angular/http";
import {CookieService, StoreService} from 'le5le-store';

import { HttpService } from '../../core/http.service';

@Injectable()
export class SitesService extends HttpService {
  public constructor(protected http: Http, protected store: StoreService) {
    super(http, store);
  }

  List(params: any): Observable<any> {
    return this.QueryString(params).Get('/api/site/list');
  }

  GetSite(params: any): Observable<any> {
    return this.QueryString(params).Get('/api/site/get');
  }

  Save(params: any): Observable<any> {
    return this.Post('/api/site/save', params);
  }
}
