import { AfterViewInit, ChangeDetectionStrategy, Component, OnInit, ViewChild } from '@angular/core';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { User } from '../../../api/models/user';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { faEdit } from '@fortawesome/free-solid-svg-icons';
import { RouterLink } from '@angular/router';
import { MatIconModule } from '@angular/material/icon';
import { NgForOf, NgOptimizedImage } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatButtonModule } from '@angular/material/button';
import { MatSelectModule } from '@angular/material/select';

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
    NgForOf
  ],
  standalone: true
})
export class UsersComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator: MatPaginator | undefined;
  public faEdit = faEdit;

  public displayedColumns: string[] = [
    'actions',
    'avatar',
    'nickname',
    'phone',
    'roles'
  ];
  public dataSource = new MatTableDataSource<User>([]);
  roles: { id: number, name: string }[] = [
    {id: 1, name: '管理员'},
    {id: 2, name: '教练'},
    {id: 3, name: '会员'}
  ];
  public selectRoleIds: number[] = [];

  constructor(private adminApi: AdminApiService) {
  }

  ngOnInit(): void {
  }

  refreshTable() {
    if (this.paginator === undefined) {
      return;
    }
    this.adminApi
      .getUsers({
        desc: true,
        keyword: '',
        limit: this.paginator.pageSize,
        page: this.paginator.pageIndex,
        roleIds: this.selectRoleIds,
      })
      .subscribe((data) => {
        this.dataSource.data = data.users;
        if (this.paginator === undefined) {
          return;
        }
        this.paginator.pageSize = data.limit;
        this.paginator.pageIndex = data.page;
        this.paginator.length = data.count;
      });
  }

  ngAfterViewInit(): void {
    console.log('after view init', this.paginator);
    this.refreshTable();

    if (this.paginator === undefined) {
      return;
    }
    this.paginator.page.subscribe(() => {
      this.refreshTable();
    });
  }
}
