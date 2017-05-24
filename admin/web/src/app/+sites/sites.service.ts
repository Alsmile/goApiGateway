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

  async SaveApi(params: any): Promise<any> {
    let ret = await this.http.Post('/api/site/api/save', params);
    if (!ret || ret.error) return {};

    return ret;
  }

  async DelApi(params: any): Promise<any> {
    let ret = await this.http.QueryString(params).Get('/api/site/api/del');
    if (!ret || ret.error) return false;

    return true;
  }

  async GetApiList(params: any): Promise<any> {
    let ret = await this.http.QueryString(params).Get('/api/site/api/list');
    if (!ret || ret.error) return [];

    return ret.list || [];
  }

  async GetApi(params: any): Promise<any> {
    let ret = await this.http.QueryString(params).Get('/api/site/api/get');
    if (!ret || ret.error) return {};

    return ret;
  }

  getMockObject(arr: any[]): any {
    if (!arr) return {};

    let mock: any = {};
    let curNode: any = mock;
    let parentNode: any = null;
    for (let item of arr) {
      if (!item.parentId) {
        curNode = mock;
        parentNode = null;
      }

      if (!item.type) {
        curNode[item.name] = item.mock;
        continue;
      }
      if (item.type.indexOf('array') === 0 && item.type !== 'array<object>') {
        if (item.mock) curNode[item.name] = JSON.parse(item.mock);
        continue;
      }

      if (item.type === 'object' || item.type === 'array<object>') {
        if (item.mock) curNode[item.name] = JSON.parse(item.mock);
        if (!item.hasChild) continue;
      }

      curNode[item.name] = item.mock;

      if (item.hasChild) {
        if (item.type === 'array<object>') {
          curNode[item.name] = [{}];
          parentNode = curNode;
          curNode = curNode[item.name][0];
        }
        else {
          curNode[item.name] = {};
          parentNode = curNode;
          curNode = curNode[item.name];
        }
      }
    }

    return mock;
  }

  getMockText(arr: any[]): string {
    if (!arr) return '';

    return JSON.stringify(this.getMockObject(arr));
  }

  strObjToArr (strObj: any): any[] {
    let ret: any[] = [];
    try {
      let obj = JSON.parse(strObj);
      ret = this.objToArr(obj);
    } catch (error) {}

    return ret;
  }

  objToArr (obj: any, parent?: any): any[] {
    let ret: any[] = [];
    for (let prop in obj) {
      let tmp: any = {
        id : new Date().getTime(),
        name : prop,
        type : typeof obj[prop],
        desc : "",
        required : "false",
        mock : obj[prop],
        level : parent? parent.level+1: 1
      };
      if (parent) tmp.parentId = parent.id;

      if (!obj[prop]) {
        tmp.type = 'string';
        ret.push(tmp);
        continue;
      }

      if (toString.apply(obj[prop]) === '[object Array]') {
        tmp.type = 'array<any>';
        if (!obj[prop][0]) continue;

        let t = typeof obj[prop][0];
        if (t === 'string') {
          tmp.type = 'array<string>';
          ret.push(tmp);
          continue;
        } else if (t === 'number') {
          tmp.type = 'array<number>';
          ret.push(tmp);
          continue;
        } else if (t === 'boolean') {
          tmp.type = 'array<boolean>';
          ret.push(tmp);
          continue;
        } else if (t === 'object') {
          tmp.type = 'array<object>';
          tmp.mock = '';
          let children = this.objToArr(obj[prop][0], tmp);
          if (children && children.length > 0) tmp.hasChild = true;
          ret.push(tmp);
          ret = ret.concat(children);
          continue;
        }
        continue;
      } else if (tmp.type !== 'object') {
        ret.push(tmp);
        continue;
      }

      tmp.mock = '';
      let children = this.objToArr(obj[prop], tmp);
      if (children && children.length > 0) tmp.hasChild = true;
      ret.push(tmp);
      ret = ret.concat(children);
    }

    return ret;
  }

  parseRequestData(arr: any, data: any) {
    if (!arr || !arr[0]) return;

    let level = 1;
    let parse = function (data: any, arr: any, level: number) {
      for (let k in data) {
        for (let item of arr) {
          if (k === item.name && level == item.level) {
            item.resValue = data[k];
            switch (item.type) {
              case 'number':
                if (!isNaN(data[k])) item.expected = 1;
                break;
              case 'string':
                if (Object.prototype.toString.call(data[k]) === "[object String]" && !!data[k]) item.expected = 1;
                break;
              case 'boolean':
                if (typeof data[k] === "boolean") item.expected = 1;
                break;
              case 'object':
                item.resValue = '';
                item.expected = 0;
                if (typeof data[k] === "object") {
                  parse(data[k], item, level+1);
                } else {
                  item.expected = -1;
                }
                break;
              default :
                if (item.type && item.type.indexOf('array')===0) {
                  if (typeof data[k] === 'object' && data[k].constructor === Array) {
                    item.expected = 1;
                    item.resValue = JSON.stringify(data[k]);
                  }
                }
                break;
            }
          }
        }
      }
    };

    parse(data, arr, level);
  }
}
