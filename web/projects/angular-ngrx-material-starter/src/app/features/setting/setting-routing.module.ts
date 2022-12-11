import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { SettingComponent } from './setting.component';

const routes: Routes = [
  {
    path: '',
    component: SettingComponent,
    data: {title: 'anms.menu.setting'},
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SettingRoutingModule {
}
