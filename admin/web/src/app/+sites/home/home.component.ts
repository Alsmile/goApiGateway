import {Component} from '@angular/core';
import {Router} from "@angular/router";

import {SitesService} from "../sites.service";

@Component({
  selector: 'sites-home',
  templateUrl: "home.component.html"
})
export class SitesHomeComponent {
  sites: any[] = [];
  pageIndex: number = 1;
  pageCount: number = 10;
  constructor(private _sitesService: SitesService, private _router: Router) {
  }

  async ngOnInit(): Promise<void> {
    this.sites = await this._sitesService.List({pageIndex: this.pageIndex, pageCount: this.pageCount});
  }
}
