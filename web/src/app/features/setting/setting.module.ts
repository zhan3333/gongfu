import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '../../shared/shared.module';
import { SettingComponent } from './setting.component';
import { SettingRoutingModule } from './setting-routing.module';
import { MatRippleModule } from '@angular/material/core';



@NgModule({
  declarations: [SettingComponent],
  imports: [
    SharedModule, SettingRoutingModule, MatRippleModule
  ]
})
export class SettingModule { }
