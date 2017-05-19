import {Component, AfterViewChecked, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import {Router, ActivatedRoute} from "@angular/router";
import {StoreService} from 'le5le-store';
import {NoticeService} from "le5le-components";

import {SitesService} from "../sites.service";

@Component({
  selector: 'sites-apis-list',
  templateUrl: "list.component.html"
})
export class SitesApisListComponent{
  loading: boolean = true;
  id: string;
  user: any;
  site: any = {};
  treeStyle: any = {};
  tree: any = {
    edited: [],
    found: [],
    showEdited: true,
    showFound: false,
    activeEdited: true,
    activeFound: false,
    selected: {}
  };
  pageIndex: number = 1;
  pageCount: number = 100;
  saving: boolean;
  constructor(private _sitesService: SitesService, private _storeService: StoreService,
              private _router: Router, private _activateRoute: ActivatedRoute) {
    this.user = _storeService.get('user');

    this.treeStyle = {
      overflow: 'auto',
      height: (document.documentElement.clientHeight - 159) + 'px'
    };
  }

  async ngOnInit(): Promise<any> {
    this.loading = true;
    this.id = this._activateRoute.snapshot.queryParams['id'];
    if (!this.id) return this.loading = false;

    this.site = await this._sitesService.GetSite({id: this.id});
    this.tree.edited = await this._sitesService.GetApiList({siteId: this.id, pageIndex: this.pageIndex, pageCount: this.pageCount});
    if (this.tree.edited.length > 0) {
      await this.onSelectEdited(this.tree.edited[0]);
    }
    this.loading = false;
  }

  async onSelectEdited(item: any): Promise<any> {
    this.tree.selected = await this._sitesService.GetApi({id: item.id});
    if (this.tree.selected.bodyParamsText)
      this.tree.selected.bodyParams = this.strObjToArr(this.tree.selected.bodyParamsText);
    if (this.tree.selected.responseParamsText)
      this.tree.selected.responseParams = this.strObjToArr(this.tree.selected.responseParamsText);
  }

  onTreeShowEdited() {
    this.tree.showEdited = !this.tree.showEdited;
    this.tree.activeEdited = true;
    this.tree.activeFound = false;
  }

  onTreeShowFound() {
    this.tree.showFound = !this.tree.showFound;
    this.tree.activeEdited = false;
    this.tree.activeFound = true;
  }

  onAdd() {
    let _noticeService: NoticeService = new NoticeService();
    _noticeService.input({
      title: '添加自定义API',
      text: '',
      placeholder: '请输入名称',
      required: true,
      callback: async (retText: any): Promise<void> => {
        let api: any = {
          owner: this.user,
          editor: this.user,
          site: this.site,
          name: retText
        };
        let ret = await this._sitesService.SaveApi(api);
        if (!ret.id) return;

        api.id = ret.id;
        api.active = true;
        this.tree.selected = api;
        this.tree.edited.push(api);
      }
    });
  }

  onSaveMock(item) {
    if (item.isEdit) return;
    this.onSaveApi();
  }

  async onSaveApi(): Promise<any> {
    this.saving = true;
    this.tree.selected.editor = this.user;
    if (this.tree.selected.contentType === 'application/json' ||
      this.tree.selected.contentType === 'multipart/form-data' ||
      this.tree.selected.contentType === 'application/x-www-form-urlencoded') {
      this.tree.selected.bodyParamsText = this.getMockText(this.tree.selected.bodyParams);
    }

    if (this.tree.selected.dataType === 'application/json' ||
      this.tree.selected.dataType === 'multipart/form-data' ||
      this.tree.selected.dataType === 'application/x-www-form-urlencoded') {
      this.tree.selected.responseParamsText = this.getMockText(this.tree.selected.responseParams);
    }

    this.tree.selected.site = this.site;
    let ret = await this._sitesService.SaveApi(this.tree.selected);
    this.saving = false;
    if (!ret.id) return;
    this.tree.selected.isEdit=false;
  }

  getMockText(arr: any[]): string {
    if (!arr) return '';

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

    return JSON.stringify(mock);
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
}
