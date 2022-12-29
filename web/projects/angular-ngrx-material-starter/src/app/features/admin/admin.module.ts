import { NgModule } from '@angular/core';
import { SharedModule } from '../../shared/shared.module';
import { UsersComponent } from './users/users.component';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { AdminRoutingModule } from './admin-routing.module';

@NgModule({
  declarations: [UsersComponent],
  imports: [
    SharedModule,
    MatTableModule,
    MatPaginatorModule,
    AdminRoutingModule,
  ]
})
export class AdminModule {
}
