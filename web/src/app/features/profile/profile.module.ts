import { NgModule } from '@angular/core';
import { CommonModule, NgOptimizedImage } from '@angular/common';
import { SharedModule } from '../../shared/shared.module';
import { ProfileRoutingModule } from './profile-routing.module';
import { ProfileComponent } from './profile.component';
import { MatLineModule } from '@angular/material/core';

@NgModule({
    imports: [
        CommonModule,
        SharedModule,
        ProfileRoutingModule,
        MatLineModule,
        NgOptimizedImage,
        ProfileComponent
    ]
})
export class ProfileModule {}
