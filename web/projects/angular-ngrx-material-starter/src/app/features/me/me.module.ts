import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MeComponent } from './me.component';
import { MeRoutingModule } from './me-routing.module';
import { SharedModule } from '../../shared/shared.module';
import { MatLegacyProgressBarModule as MatProgressBarModule } from '@angular/material/legacy-progress-bar';

@NgModule({
  declarations: [MeComponent],
  imports: [CommonModule, SharedModule, MeRoutingModule, MatProgressBarModule]
})
export class MeModule {}
