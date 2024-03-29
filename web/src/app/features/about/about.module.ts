import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { SharedModule } from '../../shared/shared.module';

import { AboutComponent } from './about/about.component';
import { AboutRoutingModule } from './about-routing.module';

@NgModule({
    imports: [CommonModule, SharedModule, AboutRoutingModule, AboutComponent]
})
export class AboutModule {}
