import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { SharedModule } from '../../shared/shared.module';

import { FeatureListComponent } from './feature-list/feature-list.component';
import { FeatureListRoutingModule } from './feature-list-routing.module';

@NgModule({
    imports: [CommonModule, SharedModule, FeatureListRoutingModule, FeatureListComponent]
})
export class FeatureListModule {}
