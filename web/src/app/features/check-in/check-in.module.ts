import { NgModule } from '@angular/core';
import { SharedModule } from '../../shared/shared.module';
import { CheckInComponent } from './check-in.component';
import { CheckInRoutingModule } from './check-in-routing.module';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { CheckInShowComponent } from './check-in-show/check-in-show.component';
import { CheckInHistoriesComponent } from './check-in-histories/check-in-histories.component';
import { MatExpansionModule } from '@angular/material/expansion';

@NgModule({
  declarations: [
    CheckInComponent,
    CheckInShowComponent,
    CheckInHistoriesComponent
  ],
  imports: [
    SharedModule,
    CheckInRoutingModule,
    MatProgressBarModule,
    MatExpansionModule
  ]
})
export class CheckInModule {}
