import {Component} from '@angular/core';
import {Router} from "@angular/router";

import {SitesHomeService} from "./home.service";

@Component({
  selector: 'sites-home',
  templateUrl: "home.component.html"
})
export class SitesHomeComponent {

  constructor(private _sitesHomeService: SitesHomeService, private _router: Router) {
  }

}
