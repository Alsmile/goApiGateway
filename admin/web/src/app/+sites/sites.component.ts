import {Component} from '@angular/core';

@Component({
  selector: 'sites-center',
  templateUrl: "sites.component.html"
})
export class SitesComponent{
  myStyle: any = {};
  constructor() {
    this.myStyle = {
      'min-height': (document.documentElement.clientHeight - 120) + 'px'
    };
  }
}
