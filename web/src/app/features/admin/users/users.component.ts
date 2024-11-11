import { AfterViewInit, ChangeDetectionStrategy, Component, OnInit, ViewChild } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { NgForOf, NgOptimizedImage } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';
import { ReactiveFormsModule } from '@angular/forms';
import { UsersTableService } from './users-table.service';

@Component({
  selector: 'anms-users',
  templateUrl: './users.component.html',
  changeDetection: ChangeDetectionStrategy.Default,
  imports: [
    MatPaginatorModule,
    MatTableModule,
    RouterLink,
    MatIconModule,
    NgOptimizedImage,
    FontAwesomeModule,
    MatButtonModule,
    MatSelectModule,
    NgForOf,
    ReactiveFormsModule
  ],
  standalone: true
})
export class UsersComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator!: MatPaginator;

  public displayedColumns: string[] = [
    'actions',
    'avatar',
    'nickname',
    'phone',
    'roles'
  ];

  roles: { id: number, name: string }[] = [
    {id: 1, name: '管理员'},
    {id: 2, name: '教练'},
    {id: 3, name: '用户'},
    {id: 4, name: '会员'}
  ];

  constructor(
    public userTableService: UsersTableService,
  ) {

  }

  ngOnInit(): void {
  }


  ngAfterViewInit(): void {
    console.log('after view init', this.paginator);
    this.userTableService.setPaginator(this.paginator)
    this.userTableService.refreshTable();
  }
}
