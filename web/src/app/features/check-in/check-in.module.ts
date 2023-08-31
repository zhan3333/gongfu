import { NgModule } from '@angular/core';
import { SharedModule } from '../../shared/shared.module';
import { CheckInComponent } from './check-in.component';
import { CheckInRoutingModule } from './check-in-routing.module';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { CheckInShowComponent } from './check-in-show/check-in-show.component';
import { CheckInHistoriesComponent } from './check-in-histories/check-in-histories.component';
import { MatExpansionModule } from '@angular/material/expansion';

@NgModule({
    imports: [
        SharedModule,
        CheckInRoutingModule,
        MatProgressBarModule,
        MatExpansionModule,
        CheckInComponent,
        CheckInShowComponent,
        CheckInHistoriesComponent
    ]
})
export class CheckInModule {}
