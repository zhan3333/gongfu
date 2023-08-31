import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { SharedModule } from '../../shared/shared.module';

import { SettingsRoutingModule } from './settings-routing.module';
import { SettingsContainerComponent } from './settings/settings-container.component';

@NgModule({
    imports: [CommonModule, SharedModule, SettingsRoutingModule, SettingsContainerComponent]
})
export class SettingsModule {}
