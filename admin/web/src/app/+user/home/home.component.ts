import {Component} from '@angular/core';
import {Router} from "@angular/router";

import {UserHomeService} from "./home.service";

@Component({
  selector: 'user-home',
  templateUrl: "home.component.html"
})
export class UserHomeComponent {

  constructor(private _userHomeService: UserHomeService, private _router: Router) {
  }

}
