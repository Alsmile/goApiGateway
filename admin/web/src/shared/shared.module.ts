import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { Le5leComponentsModule } from "le5le-components";

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    Le5leComponentsModule,
  ],
  declarations: [ ],
  exports: [
    CommonModule,
    FormsModule,
    Le5leComponentsModule,
  ],
  providers: [ ]
})
export class SharedModule {
}
