import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { MeComponent } from './me.component';

const routes: Routes = [
  {
    path: '',
    component: MeComponent,
    data: {title: 'anms.menu.me'}
  }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MeRoutingModule {
}
