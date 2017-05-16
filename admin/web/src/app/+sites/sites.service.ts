import { Injectable } from '@angular/core';
import { StoreService } from 'le5le-store';

import { HttpService } from '../../core/http.service';

@Injectable()
export class SitesService {
  public constructor(protected http: HttpService, protected store: StoreService) {
  }

  async List(params: any): Promise<any> {
    let ret = await this.http.QueryString(params).Get('/api/site/list');
    if (!ret || !ret.list) return [];
    return ret.list;
  }

  async GetSite(params: any): Promise<any> {
    let ret = await this.http.QueryString(params).Get('/api/site/get');
    if (!ret || ret.error) return {https: '', notFound: {code:404}, statics: [], proxies: []};

    return ret;
  }

  async Save(params: any): Promise<any> {
    let ret = await this.http.Post('/api/site/save', params);
    if (!ret || ret.error) return false;

    return true;
  }
}
