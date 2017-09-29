import {Component, Input} from '@angular/core';
import {HttpService} from '../../../core/http.service';
import {SitesService} from '../sites.service';

@Component({
  selector: 'api-edit',
  templateUrl: 'api.edit.component.html'
})
export class ApiEditComponent {
  @Input() api: any = {};
  requestData: any;
  constructor(protected http: HttpService, private _sitesService: SitesService) {
  }

  ngOnChanges() {
    this.api.method = this.api.method || 'GET';
    this.api.contentType = this.api.contentType || 'application/json';
    this.api.dataType = this.api.dataType || 'application/json';
  }

  onAddHeader() {
    if (!this.api.headers) this.api.headers = [];
    this.api.headers.push({
      name: '',
      required: 'false',
      desc: ''
    });
  }
  onAddQuery() {
    if (!this.api.queryParams) this.api.queryParams = [];
    this.api.queryParams.push({
      name: '',
      required: 'false',
      desc: ''
    });
  }

  canShow (arr: any, i: number) {
    if (i < 1) return true;

    let show: boolean = true;
    let level = arr[i].level;
    for (let j = i - 1; j > -1; --j) {
      if (show && arr[j].level < level) {
        level = arr[j].level;
        show = arr[j].hide ? false : true;
      }
    }

    return show;
  }

  canDel(arrs: any[], item: any) {
    let isFind: boolean;
    for (let val of arrs) {
      if (val.parentId === item.id) isFind = true;
    }
    return !isFind;
  }

  getMarginLeft(item: any) {
    if (!item || item.level < 2) return '.2rem';

    return 0.2 * item.level + 'rem';
  }

  onAddBody(parentItem: any, pos?: number) {
    if (!this.api.bodyParams) this.api.bodyParams = [];
    let newItem: any = {
      id: new Date().getTime(),
      name: '',
      type: 'string',
      required: 'false',
      desc: '',
      mock: '',
      level: 1
    };

    if (parentItem) {
      parentItem.hasChild = true;
      newItem.parentId = parentItem.id;
      newItem.level += parentItem.level;
      let i:number = pos+1;
      for (; i < this.api.bodyParams.length - 1; ++i) {
        if (this.api.bodyParams[i].level < newItem.level) break;
      }
      this.api.bodyParams.splice(i, 0, newItem);
    }
    else {
      this.api.bodyParams.push(newItem);
    }
  }

  onAddResponse(parentItem: any, pos?: number) {
    if (!this.api.responseParams) this.api.responseParams = [];
    let newItem: any = {
      id: new Date().getTime(),
      name: '',
      type: 'string',
      required: 'false',
      desc: '',
      mock: '',
      level: 1
    };

    if (parentItem) {
      parentItem.hasChild = true;
      newItem.parentId = parentItem.id;
      newItem.level += parentItem.level;
      let i:number = pos+1;
      for (; i < this.api.responseParams.length; ++i) {
        if (this.api.responseParams[i].level < newItem.level) break;
      }

      this.api.responseParams.splice(i, 0, newItem);
    }
    else {
      this.api.responseParams.push(newItem);
    }
  }

  requestRet: any = '';
  async onRequestApi(): Promise<any> {
    let headers: any = {};
    for (let item of this.api.headers) {
      headers[item.name] = item.mock;
    }

    let queryParams: any = {
      host: this.api.site.dstUrl,
      url: this.api.url,
      dataType: this.api.dataType
    };
    for (let item of this.api.queryParams) {
      queryParams[item.name] = item.mock;
    }

    let bodyParams: any = this.api.bodyParamsText;
    if (this.api.contentType === 'application/json' ||
      this.api.contentType === 'multipart/form-data' ||
      this.api.contentType === 'application/x-www-form-urlencoded') {
      this._sitesService.getMockObject(this.api.bodyParams);
    }
    if (this.api.method === 'GET') {
      this.requestRet =
        await this.http.QueryString(queryParams).Get('/api/test', { headers: headers });
    } else if (this.api.method === 'POST') {
      this.requestRet =
        await this.http.QueryString(queryParams).Post('/api/test', bodyParams, { headers: headers });
    } else if (this.api.method === 'PUT') {
      this.requestRet =
        await this.http.QueryString(queryParams).Put('/api/test', bodyParams, { headers: headers });
    } else if (this.api.method === 'DELETE') {
      this.requestRet =
        await this.http.QueryString(queryParams).Delete('/api/test', { headers: headers });
    }

    this._sitesService.parseRequestData(this.api.responseParams, this.requestRet);
  }

  extractUrlParams(str: string) {
    let ret: any = {
      urlReg: '',
      urlParams: [],
      urlParamSrcs: []
    };

    if (!str) return ret;

    ret.urlParamSrcs = str.match(/\/<.*?>/g);
    if (!ret.urlParamSrcs) return ret;

    for (let item of ret.urlParamSrcs) {
      let s = item.replace(/\/<(.+?)\s*>/, '$1').split('=')[0];
      if (!s || s === item) continue;
      let isFind = false;
      for (let i of ret.urlParams) {
        if (i === s) {
          isFind = true;
          break;
        }
      }

      if (!isFind) {
        let newItem: any = { name: s };
        for (let i of this.api.urlParams) {
          if (i.name === s) {
            newItem.desc = i.desc;
            break;
          }
        }
        ret.urlParams.push(newItem);
      }

      let arr = item.split('=');
      if (arr[1]) str = str.replace(item, '/' + arr[1].substr(0, arr[1].length - 1));
      else str = str.replace(item, '/(.+)');
    }
    ret.urlReg = '^' + str.replace(/\//g, '\/') + '$';

    return ret;
  }

  onChangedUrl() {
    if (!this.api.shortUrl) {
      this.api.urlReg = '';
      this.api.urlParams = [];
      return;
    }

    let index = this.api.shortUrl.indexOf('?');
    if (index === 0) {
      this.api.shortUrl = '';
      this.api.urlReg = '';
      this.api.urlParams = [];
      return;
    }
    else if (index > 0) {
      this.api.shortUrl = this.api.shortUrl.substr(0, index);
    }

    let urlParams = this.extractUrlParams(this.api.shortUrl);
    this.api.urlParams = urlParams.urlParams;
    this.api.urlReg = urlParams.urlReg;
  }
}
