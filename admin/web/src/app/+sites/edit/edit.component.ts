import {Component, AfterViewChecked, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
import {Router, ActivatedRoute} from "@angular/router";
import {StoreService} from 'le5le-store';
import {NoticeService} from "le5le-components";

import {SitesService} from "../sites.service";

@Component({
  selector: 'site-edit',
  templateUrl: "edit.component.html"
})
export class SitesEditComponent implements AfterViewChecked {
  loading: boolean = true;
  id: string;
  user: any;
  site: any = {https: '', notFound: {code:404}, isSubdomain: true};
  saving: boolean;
  formErrors: any = {};
  domain: string = document.domain.replace('admin', '');
  @ViewChild('myForm') currentForm: NgForm;
  constructor(private _sitesService: SitesService, private _storeService: StoreService,
              private _router: Router, private _activateRoute: ActivatedRoute) {
    this.user = _storeService.get('user');
    this.site.owner = this.site.editor = this.user;
  }

  async ngOnInit(): Promise<any> {
    this.loading = true;
    this.id = this._activateRoute.snapshot.queryParams['id'];
    if (!this.id) return this.loading = false;

    this.site = await this._sitesService.GetSite({id: this.id});
    this.loading = false;
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
    for (let item in this.formErrors) this.formErrors[item] = false;
    for (const field in form.controls) {
      const control = form.get(field);

      if (control && (dirty || control.dirty) && !control.valid) {
        this.formErrors[field] = true;
      }
    }
  }

  async onSubmit(): Promise<void> {
    this.onValueChanged(true);
    if (this.currentForm.form.invalid) return;

    this.saving = true;
    this.site.editor = this.user;
    if (!this.site.isCustomDomain) {
      if (this.site.subdomain) this.site.apiDomain = this.site.subdomain + this.domain;
      else this.site.apiDomain = this.domain.substring(1);
    }
    let ret = await this._sitesService.Save(this.site);
    if (ret) {
      let _noticeService: NoticeService = new NoticeService();
      _noticeService.notice({body: '保存成功！', theme: 'success'});

      this._router.navigateByUrl('/sites/home');
    }
    this.saving = false;
  }
}
