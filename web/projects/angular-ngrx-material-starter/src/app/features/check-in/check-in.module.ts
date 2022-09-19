import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '../../shared/shared.module';
import { CheckInComponent } from './check-in.component';
import { CheckInRoutingModule } from './check-in-routing.module';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { CheckInShowComponent } from './check-in-show/check-in-show.component';
import { CheckInTopComponent } from './check-in-top/check-in-top.component';


@NgModule({
  declarations: [CheckInComponent, CheckInShowComponent, CheckInTopComponent],
  imports: [
    CommonModule, SharedModule, CheckInRoutingModule, MatProgressBarModule
  ]
})
export class CheckInModule {
}
