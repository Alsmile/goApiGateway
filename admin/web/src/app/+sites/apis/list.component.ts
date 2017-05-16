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
  constructor(private _sitesService: SitesService, private _storeService: StoreService,
              private _router: Router, private _activateRoute: ActivatedRoute) {
    this.user = _storeService.get('user');

    this.treeStyle = {
      'height': (document.documentElement.clientHeight - 159) + 'px'
    };
  }

  async ngOnInit(): Promise<any> {
    this.loading = true;
    this.id = this._activateRoute.snapshot.queryParams['id'];
    if (!this.id) return this.loading = false;

    this.site = await this._sitesService.GetSite({id: this.id});
    this.loading = false;
  }

  onTreeShowEdited(item: any) {
    this.tree.showEdited = !this.tree.showEdited;
    this.tree.activeEdited = true;
    this.tree.activeFound = false;
  }

  onTreeShowFound(item: any) {
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

        let ret = await this._sitesService.SaveApi({owner: this.user, siteId: this.site.id, name: retText});
        if (!ret.id) return;

        this.tree.edited.push({
          id: ret.id,
          name: retText,
          active: true
        });
      }
    });
  }
}
