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
  foundApis: any[] = [];
  foundIndex: number = 1;
  foundCount: number = 10;
  loadingFoundApis: boolean = false;
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
    this.tree.edited = await this._sitesService.GetApiList({
      siteId: this.id,
      auto: 'false',
      field: 1,
      pageIndex: this.pageIndex,
      pageCount: this.pageCount
    });
    if (this.tree.edited.length > 0) {
      await this.onSelectEdited(this.tree.edited[0]);
    }
    this.loading = false;
  }

  async onSelectEdited(item: any): Promise<any> {
    this.tree.activeEdited = true;
    this.tree.activeFound = false;
    this.tree.selected = await this._sitesService.GetApi({id: item.id});
    if ((!this.tree.selected.bodyParams || this.tree.selected.bodyParams.length < 1) && this.tree.selected.bodyParamsText)
      this.tree.selected.bodyParams = this._sitesService.strObjToArr(this.tree.selected.bodyParamsText);
    if ((!this.tree.selected.responseParams || this.tree.selected.responseParams.length < 1) && this.tree.selected.responseParamsText)
      this.tree.selected.responseParams = this._sitesService.strObjToArr(this.tree.selected.responseParamsText);
  }

  onTreeShowEdited() {
    this.tree.showEdited = !this.tree.showEdited;
    this.tree.activeEdited = true;
    this.tree.activeFound = false;
  }

  onTreeShowFound() {
    this.tree.activeEdited = false;
    this.tree.activeFound = true;
    this.tree.selected = {};

    this.foundIndex = 1;
    this.getFoundApis();
  }

  async getFoundApis(): Promise<any> {
    this.loadingFoundApis = true;
    this.foundApis = await this._sitesService.GetApiList({
      siteId: this.id,
      auto: 'true',
      pageIndex: this.foundIndex,
      pageCount: this.foundCount
    });
    this.loadingFoundApis = false;
  }

  onLastPage() {
    --this.foundIndex;
    this.getFoundApis();
  }

  onNextPage() {
    ++this.foundIndex;
    this.getFoundApis();
  }

  onEditFound(item: any){
    item.isEdit = true;
    this.tree.selected = item;
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

  onDelApi() {
    let _noticeService: NoticeService = new NoticeService();
    _noticeService.dialog({
      title: '确认',
      theme: '',
      body: '确认删除此api？',
      callback: async (ret:boolean): Promise<any> =>{
        if (!await this._sitesService.DelApi({id: this.tree.selected.id})) return;

        for (let i=0; i < this.tree.edited.length; i++) {
          if (this.tree.selected.id === this.tree.edited[i].id) {
            this.tree.edited.splice(i,1);
            this.tree.selected = {};
            if (this.tree.edited.length > 0) {
              this.onSelectEdited(this.tree.edited[0]);
            }
          }
        }
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
      this.tree.selected.bodyParamsText = this._sitesService.getMockText(this.tree.selected.bodyParams);
    }

    if (this.tree.selected.dataType === 'application/json' ||
      this.tree.selected.dataType === 'multipart/form-data' ||
      this.tree.selected.dataType === 'application/x-www-form-urlencoded') {
      this.tree.selected.responseParamsText = this._sitesService.getMockText(this.tree.selected.responseParams);
    }

    this.tree.selected.site = this.site;
    this.tree.selected.autoReg = false;
    this.tree.selected.url = this.site.group + this.tree.selected.shortUrl;
    let ret = await this._sitesService.SaveApi(this.tree.selected);
    if (this.tree.activeFound) {
      this.tree.edited.push(this.tree.selected);
    }
    this.saving = false;
    if (!ret.id) return;
    this.tree.selected.isEdit=false;
  }

  onCancelEdit() {
    if (this.tree.activeFound) {
      this.tree.selected = {};
    } else {
      this.onSelectEdited(this.tree.selected);
    }
  }
}
