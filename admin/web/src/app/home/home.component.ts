import {Component} from '@angular/core';
import {Router} from "@angular/router";
import {StoreService} from 'le5le-store';

import {HomeService} from "./home.service";

@Component({
  selector: 'home',
  templateUrl: "home.component.html",
  providers: [HomeService]
})
export class HomeComponent {
  constructor(private _homeService: HomeService, private _router: Router, private _storeService: StoreService) {
  }

  onNeedActive () {
    this._storeService.set('needActive', 1);
  }

}
