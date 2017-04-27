import {Injectable} from '@angular/core';
import {Http, Headers, Response} from "@angular/http";
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';

import {NoticeService} from "le5le-components";
import {CookieService, StoreService} from 'le5le-store';

@Injectable()
export class HttpService {
  baseUrl: string = '';
  queryParams: string = '';
  headers: Headers = new Headers();
  public constructor(protected http: Http, protected store: StoreService) {
    this.headers.append('Content-Type', 'application/json');
  }

  private getToken (): string {
    let remember:any = localStorage.getItem("rememberMe");
    if (remember) {
      return localStorage.getItem("token");
    } else {
      return CookieService.get("token");
    }
  }

  private delToken () {
    this.store.set('auth', false);
  }

  SetBaseUrl(url: string): HttpService {
    this.baseUrl = url;
    return this;
  }

  QueryString(obj: any): HttpService {
    this.queryParams = '?' +
      Object.keys(obj).map(function(key) {
        if (!obj[key]) return '';

        if (obj[key] instanceof Array || Object.prototype.toString.call((obj[key])) == '[object Array]') {
          return obj[key].map(function(item: string) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(item);
          }).join('&');
        } else {
          return encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]);
        }
      }).join('&');

    return this;
  }

  AppendHeader(name: string, value: string): HttpService {
    this.headers.set(name, value);
    return this;
  }

  setHeaders(options?: any) {
    this.headers.set('Authorization', 'Bearer ' + this.getToken());

    if (!options || !options.headers) return;

    Object.keys(options.headers).map(function(key) {
      if (options.headers[key]) this.headers.set(key, options.headers[key]);
    })
  }

  Get(url: string, options?: any): Observable<any> {
    this.setHeaders(options);
    url += this.queryParams;
    this.queryParams = '';
    return this.http.get(this.baseUrl + url, {headers: this.headers}).map(this.extractData)
      .filter(data => !data.error).catch(this.handleError);
  }

  Delete(url: string, options?: any): Observable<any> {
    this.setHeaders(options);
    url += this.queryParams;
    this.queryParams = '';
    return this.http.delete(this.baseUrl + url, {headers: this.headers}).map(this.extractData)
      .filter(data => !data.error).catch(this.handleError);
  }

  Post(url: string, body: any, options?: any): Observable<any> {
    let strBody: string;
    if(typeof body === 'string') {
      strBody = body;
    } else {
      strBody = JSON.stringify(body);
    }
    this.setHeaders(options);
    url += this.queryParams;
    this.queryParams = '';
    return this.http.post(this.baseUrl + url, strBody, {headers: this.headers}).map(this.extractData)
      .filter(data => !data.error).catch(this.handleError);
  }

  Put(url: string, body: any, options?: any): Observable<any> {
    let strBody: string;
    if(typeof body === 'string') {
      strBody = body;
    } else {
      strBody = JSON.stringify(body);
    }
    this.setHeaders(options);
    url += this.queryParams;
    this.queryParams = '';
    return this.http.put(this.baseUrl + url, strBody, {headers: this.headers}).map(this.extractData)
      .filter(data => !data.error).catch(this.handleError);
  }

  private extractData(res: Response) {
    if (res.text() ===  'null') return {};

    let body = res.json();
    if (body.error) {
      let _noticeService: NoticeService = new NoticeService();
      _noticeService.notice({body: body.error, theme: 'error', timeout: 5000});
    }
    return body || {};
  }

  private handleError (error: any) {
    let errMsg = error.message;
    if (error.status == 401) {
      this.delToken();
      errMsg = 'Authorization error';
    }
    console.error(errMsg);
    return Observable.throw(errMsg);
  }
}
