import {Component, AfterViewChecked, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import {Router, ActivatedRoute} from "@angular/router";

import {SitesEditService} from "./edit.service";

@Component({
  selector: 'site-edit',
  templateUrl: "edit.component.html"
})
export class SitesEditComponent implements AfterViewChecked {
  id: string;
  site: any = {https: '', statics: [], proxies: []};
  staticUrl : string;
  staticPath: string;
  proxyUrl: string;
  proxyPath: string;
  saving: boolean;
  formErrors: any = {};
  @ViewChild('myForm') currentForm: NgForm;
  constructor(private _sitesEditService: SitesEditService, private _router: Router, private _activateRoute: ActivatedRoute) {
  }

  ngOnInit() {
    this.id = this._activateRoute.snapshot.queryParams['id'];
    if (!this.id) return;

    this._sitesEditService.GetSite({id: this.id}).subscribe(
      ret => {
        this.site = ret;
      },
      err => console.error(err)
    );
  }

  ngAfterViewChecked() {
    this.formChanged();
  }

  isInitCurrentForm: boolean = false;
  formChanged() {
    if (this.currentForm && !this.isInitCurrentForm) {
      this.isInitCurrentForm = true;
      this.currentForm.valueChanges.subscribe(data => this.onValueChanged());
    }
  }

  onValueChanged(dirty?: boolean) {
    if (!this.currentForm) { return; }
    const form = this.currentForm.form;
    for (const field in form.controls) {
      this.formErrors[field] = false;
      const control = form.get(field);

      if (control && (dirty || control.dirty) && !control.valid) {
        this.formErrors[field] = true;
      }
    }
  }

  onAddStatic () {
    if (!this.staticUrl || !this.staticPath) return;

    this.site.statics.push({
      url: this.staticUrl,
      path: this.staticPath
    });
    this.staticUrl = '';
    this.staticPath = '';
  }

  onAddProxy () {
    if (!this.proxyUrl || !this.proxyPath) return;

    this.site.proxies.push({
      url: this.proxyUrl,
      path: this.proxyPath
    });

    this.proxyUrl = '';
    this.proxyPath = '';
  }

  onSubmit () {
    this.onValueChanged(true);
    if (this.currentForm.form.invalid) return;

    this.onAddStatic();
    this.onAddProxy();

    this.saving = true;
    this._sitesEditService.Save(this.site).subscribe(
      ret => {

      },
      err => console.error(err),
      () => this.saving = false
    );
  }
}
