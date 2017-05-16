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

  ngOnInit() {
    this._sitesService.List({pageIndex: this.pageIndex, pageCount: this.pageCount}).subscribe(
      ret => {
        this.sites = ret.list;
      },
      err => console.error(err)
    );
  }
}
