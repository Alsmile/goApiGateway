import {Component, Input} from '@angular/core';

@Component({
  selector: 'api-edit',
  templateUrl: 'api.edit.component.html'
})
export class ApiEditComponent {
  @Input() api: any = {};
  requestData: any;
  constructor() {
  }

  ngOnInit() {
    this.api.method = this.api.method || 'GET';
    this.api.contentType = this.api.contentType || 'application/json';
    this.api.dataType = this.api.dataType || 'application/json';
  }

  onAddHeader() {
    if (!this.api.headers) this.api.headers = [];
    this.api.headers.push({
      name: '',
      required: 'false',
      desc: ''
    });
  }
  onAddQuery() {
    if (!this.api.queryParams) this.api.queryParams = [];
    this.api.queryParams.push({
      name: '',
      required: 'false',
      desc: ''
    });
  }

  canShow (arr: any, i: number) {
    if (i < 1) return true;

    let show: boolean = true;
    let level = arr[i].level;
    for (let j=i-1; j > -1; --j) {
      if (show && arr[j].level < level) {
        level = arr[j].level;
        show = arr[j].hide?false: true;
      }
    }

    return show;
  }

  canDel(arrs: any[], item: any) {
    let isFind: boolean;
    for (let val of arrs) {
      if (val.parentId === item.id) isFind = true;
    }
    return !isFind;
  }

  getMarginLeft(item: any) {
    if (!item || item.level < 2) return '.2rem';

    return 0.2 * item.level + 'rem';
  }

  onAddBody(parentItem: any, pos?: number) {
    if (!this.api.bodyParams) this.api.bodyParams = [];
    let newItem: any = {
      id: new Date().getTime(),
      name: '',
      type: 'string',
      required: 'false',
      desc: '',
      mock: '',
      level: 1
    };

    if (parentItem) {
      newItem.parentId = parentItem.id;
      newItem.level += parentItem.level;
      let i:number = pos+1;
      for (; i < this.api.bodyParams.length - 1; ++i) {
        if (this.api.bodyParams[i].level < newItem.level) break;
      }
      this.api.bodyParams.splice(i, 0, newItem);
    }
    else {
      this.api.bodyParams.push(newItem);
    }
  }

  onAddResponse(parentItem: any, pos?: number) {
    if (!this.api.responseParams) this.api.responseParams = [];
    let newItem: any = {
      id: new Date().getTime(),
      name: '',
      type: 'string',
      required: 'false',
      desc: '',
      mock: '',
      level: 1
    };

    if (parentItem) {
      newItem.parentId = parentItem.id;
      newItem.level += parentItem.level;
      let i:number = pos+1;
      for (; i < this.api.responseParams.length - 1; ++i) {
        if (this.api.responseParams[i].level < newItem.level) break;
      }
      this.api.responseParams.splice(i, 0, newItem);
    }
    else {
      this.api.responseParams.push(newItem);
    }
  }

  onRequestApi() {

  }
}
