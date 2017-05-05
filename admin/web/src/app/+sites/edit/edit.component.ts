import {Component, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import {Router, ActivatedRoute} from "@angular/router";

import {SitesEditService} from "./edit.service";

@Component({
  selector: 'site-edit',
  templateUrl: "edit.component.html"
})
export class SitesEditComponent {
  id: string;
  site: any = {https: '', statics: [], proxies: []};
  staticUrl : string;
  staticPath: string;
  proxyUrl: string;
  proxyPath: string;
  saving: boolean;
  @ViewChild('myForm') currentForm: NgForm;
  constructor(private _sitesEditService: SitesEditService, private _router: Router, private _activateRoute: ActivatedRoute) {
  }

  ngOnInit() {
    this.id = this._activateRoute.snapshot.params['id'];
  }

  onAddStatic () {
    this.site.statics.push({
      url: this.staticUrl,
      path: this.staticPath
    });

    this.staticUrl = '';
    this.staticPath = '';
  }

  onAddProxy () {
    this.site.proxies.push({
      url: this.proxyUrl,
      path: this.proxyPath
    });

    this.proxyUrl = '';
    this.proxyPath = '';
  }

  onSubmit () {
    if (this.currentForm.form.invalid) return;

    this.saving = true;
  }
}
