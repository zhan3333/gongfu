import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from '../../shared/shared.module';
import { LoginComponent } from './login.component';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { LoginRoutingModule } from './login-routing.module';

@NgModule({
  declarations: [LoginComponent],
  imports: [
    CommonModule,
    SharedModule,
    MatProgressBarModule,
    LoginRoutingModule
  ]
})
export class LoginModule {}
