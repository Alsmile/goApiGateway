import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { Le5leComponentsModule } from "le5le-components";

import {ProxyKeyValidator} from './directives/proxyKey.directive';

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    Le5leComponentsModule,
  ],
  declarations: [
    ProxyKeyValidator
  ],
  exports: [
    CommonModule,
    FormsModule,
    Le5leComponentsModule,
    ProxyKeyValidator
  ],
  providers: []
})
export class SharedModule {
}
