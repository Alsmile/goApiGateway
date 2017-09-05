import {Injectable} from '@angular/core';
import {Http, Headers, Response, ResponseContentType} from '@angular/http';
import {Observable} from 'rxjs/Observable';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/filter';
import * as FileSaver from 'file-saver';

import {NoticeService} from 'le5le-components';
import {CookieService, StoreService} from 'le5le-store';

@Injectable()
export class HttpService {
  baseUrl: string = '';
  queryParams: string = '';
  public constructor(protected http: Http, protected store: StoreService) {
  }

  private getToken(): string {
    let remember: any = localStorage.getItem('rememberMe');
    if (remember) {
      return localStorage.getItem((<any> window).token);
    } else {
      return CookieService.get((<any> window).token);
    }
  }

  private setToken(token: string) {
    let domains = document.domain.split('.');
    let strDomain = '';
    for (let i=domains.length-1; i>0 && i>domains.length-4; --i) {
      strDomain = domains[i] + '.' + strDomain;
    }
    strDomain = strDomain.substr(0, strDomain.length-1);

    let remember: any = localStorage.getItem('rememberMe');
    if (remember) {
      localStorage.setItem((<any> window).token, token);
    } else {
      CookieService.set((<any> window).token, token, {domain: strDomain});
    }
  }

  private delToken() {
    this.store.set('auth', false);
  }

  SetBaseUrl(url: string): HttpService {
    this.baseUrl = url;
    return this;
  }

  QueryString(obj: any): HttpService {
    this.queryParams = '?' +
      Object.keys(obj).map(function (key) {
        if (obj[key] === undefined || obj[key] === null || obj[key] === '') return '';

        if (obj[key] instanceof Array || Object.prototype.toString.call((obj[key])) == '[object Array]') {
          return obj[key].map(function (item: string) {
            return encodeURIComponent(key) + '=' + encodeURIComponent(item);
          }).join('&');
        } else {
          return encodeURIComponent(key) + '=' + encodeURIComponent(obj[key]);
        }
      }).join('&');

    return this;
  }

  private setHeaders(options?: any): Headers {
    let headers: Headers = new Headers();
    headers.set('Content-Type', 'application/json');
    headers.set('Authorization', this.getToken());

    if (!options || !options.headers) return headers;

    Object.keys(options.headers).map((key: string) => {
      if (options.headers[key]) headers.set(key, options.headers[key]);
    });

    return headers;
  }

  async Get(url: string, options?: any): Promise<any> {
    url += this.queryParams;
    this.queryParams = '';
    try {
      let response = await this.http
        .get(this.baseUrl + url, {headers: this.setHeaders(options)})
        .toPromise();
      return this.extractData(response);
    } catch (error) {
      return this.handleError(error);
    }
  }

  async Delete(url: string, options?: any): Promise<any> {
    url += this.queryParams;
    this.queryParams = '';

    try {
      let response = await this.http
        .delete(this.baseUrl + url, {headers: this.setHeaders(options)})
        .toPromise();
      return this.extractData(response);
    } catch (error) {
      return this.handleError(error);
    }
  }

  async Post(url: string, body: any, options?: any): Promise<any> {
    let strBody: string;
    if (typeof body === 'string') {
      strBody = body;
    } else {
      strBody = JSON.stringify(body);
    }
    url += this.queryParams;
    this.queryParams = '';

    try {
      let response = await this.http
        .post(this.baseUrl + url, strBody, {headers: this.setHeaders(options)})
        .toPromise();
      return this.extractData(response);
    } catch (error) {
      return this.handleError(error);
    }
  }

  async PostForm(url: string, body: FormData, options?: any): Promise<any> {
    url += this.queryParams;
    this.queryParams = '';

    let headers = this.setHeaders(options);
    headers.delete('Content-Type');
    try {
      let response = await this.http
        .post(this.baseUrl + url, body, {headers: headers})
        .toPromise();
      return this.extractData(response);
    } catch (error) {
      return this.handleError(error);
    }
  }

  async Put(url: string, body: any, options?: any): Promise<any> {
    let strBody: string;
    if (typeof body === 'string') {
      strBody = body;
    } else {
      strBody = JSON.stringify(body);
    }

    url += this.queryParams;
    this.queryParams = '';
    try {
      let response = await this.http
        .put(this.baseUrl + url, strBody, {headers: this.setHeaders(options)})
        .toPromise();
      return this.extractData(response);
    } catch (error) {
      return this.handleError(error);
    }
  }

  DownloadFile(url: string, fileName: string, options?: any) {
    url += this.queryParams;
    this.queryParams = '';
    let sub = this.http.get(this.baseUrl + url, {
      headers: this.setHeaders(options),
      responseType: ResponseContentType.Blob
    }).subscribe(
      (res: Response) => {
        FileSaver.saveAs(res.blob(), fileName);
      },
      err => console.error(err),
      () => { sub.unsubscribe(); }
    );
  }

  private extractData(res: Response): any {
    if (res.headers.get((<any> window).token)) this.setToken(res.headers.get((<any> window).token));

    if (!res || !res.text() || res.text() === 'null') return null;

    let body = res.json();
    if (body.error) {
      let _noticeService: NoticeService = new NoticeService();
      _noticeService.notice({body: body.error, theme: 'error', timeout: 5000});
    } else if (body.code && body.code != 0) {
      body.error = body.message;
      let _noticeService: NoticeService = new NoticeService();
      _noticeService.notice({body: body.message, theme: 'error', timeout: 5000});
    }

    return body;
  }

  private handleError(error: any): any {
    if (!error) error = {error: '未知错误'};
    if (error.status == 401) {
      let ret = JSON.parse(error._body);
      this.store.set('loginUrl', ret.loginUrl || '');
      this.delToken();
      error.error = 'Authorization error';
    } else if (error.status != 404) {
      console.error(error);
      error = {error: error};
    }
    return error;
  }
}
